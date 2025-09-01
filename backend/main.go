package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

const Version = "1.0.2"

// Metrics
var (
	tasksCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "todolist_tasks_created_total",
		Help: "Total number of tasks created",
	})

	tasksCompleted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "todolist_tasks_completed_total",
		Help: "Total number of tasks completed",
	})

	tasksDeleted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "todolist_tasks_deleted_total",
		Help: "Total number of tasks deleted",
	})

	tasksTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "todolist_tasks_total",
		Help: "Current total number of tasks",
	})

	activeTasks = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "todolist_tasks_active",
		Help: "Current number of active (not completed) tasks",
	})

	completedTasks = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "todolist_tasks_completed",
		Help: "Current number of completed tasks",
	})

	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "todolist_http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "endpoint", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "todolist_http_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: []float64{0.1, 0.3, 0.5, 1, 2, 5},
	}, []string{"method", "endpoint"})
)

// Middleware for HTTP requests
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(ww, r)
		
		duration := time.Since(start).Seconds()
		endpoint := r.URL.Path
		
		httpRequestsTotal.WithLabelValues(r.Method, endpoint, strconv.Itoa(ww.statusCode)).Inc()
		httpRequestDuration.WithLabelValues(r.Method, endpoint).Observe(duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

type Task struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Done      bool       `json:"done"`
	CreatedAt *time.Time `json:"created_at"`
}

func AddTask(db *sql.DB, title string) (*Task, error) {
	var task Task
	err := db.QueryRow(
		"INSERT INTO tasks (title) VALUES ($1) RETURNING id, title, done, created_at", title,
	).Scan(&task.ID, &task.Title, &task.Done, &task.CreatedAt)
	if err == nil {
		tasksCreated.Inc()
		updateTaskMetrics(db)
	}
	return &task, err
}

func GetAllTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query(`SELECT id, title, done, created_at FROM tasks ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Done, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func CompleteTask(db *sql.DB, id int) error {
	_, err := db.Exec("UPDATE tasks SET done = TRUE WHERE id = $1", id)
	if err == nil {
		tasksCompleted.Inc()
		updateTaskMetrics(db)
	}
	return err
}

func DeleteTask(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err == nil {
		tasksDeleted.Inc()
		updateTaskMetrics(db)
	}
	return err
}

// Upgrade metrics (amount of tasks)
func updateTaskMetrics(db *sql.DB) {
	var total, active, completed int
	
	err := db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&total)
	if err != nil {
		log.Printf("Error getting total tasks count: %v", err)
		return
	}
	
	err = db.QueryRow("SELECT COUNT(*) FROM tasks WHERE done = FALSE").Scan(&active)
	if err != nil {
		log.Printf("Error getting active tasks count: %v", err)
		return
	}
	
	err = db.QueryRow("SELECT COUNT(*) FROM tasks WHERE done = TRUE").Scan(&completed)
	if err != nil {
		log.Printf("Error getting completed tasks count: %v", err)
		return
	}
	
	tasksTotal.Set(float64(total))
	activeTasks.Set(float64(active))
	completedTasks.Set(float64(completed))
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("Failed to create driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("Failed to create migrate instance: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Failed to apply migrations: %w", err)
	}
	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Connection to PostgreSQL
	connStr := fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")

	// Migration func
	if err := runMigrations(db); err != nil {
		log.Fatal("Migration error:", err)
	}
	fmt.Println("Migration applied successfully")

	// Update metrics when app starts
	updateTaskMetrics(db)

	// Router for metrics
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())

	// Run server of metrics on port 9090
	go func() {
		fmt.Println("Metrics server running on http://localhost:9090")
		log.Fatal(http.ListenAndServe(":9090", metricsMux))
	}()

	// Create Router
	mux := http.NewServeMux()

	// Set HTTP Router
	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := GetAllTasks(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		if err := json.NewEncoder(w).Encode(tasks); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/api/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
			return
		}

		title := r.FormValue("title")
		if title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}

		task, err := AddTask(db, title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)
	})

	mux.HandleFunc("/api/done", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := r.FormValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		if err = CompleteTask(db, id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		healthStatus := map[string]interface{}{
			"status":    "ok",
			"version":   Version,
			"service":   "todolist-backend",
			"timestamp": time.Now().Format(time.RFC3339),
			"components": map[string]string{
				"database": "ok",
			},
		}

		err := db.Ping()
		if err != nil {
			healthStatus["components"].(map[string]string)["database"] = "error"
			healthStatus["status"] = "degraded"
			log.Printf("Database health check failed: %v", err)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(healthStatus)
	})

	mux.HandleFunc("/api/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := r.FormValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}
		if err = DeleteTask(db, id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// Add endpoint for metrics
	mux.HandleFunc("/api/metrics/health", func(w http.ResponseWriter, r *http.Request) {
		metricsHealth := map[string]interface{}{
			"status":    "ok",
			"metrics":   "enabled",
			"timestamp": time.Now().Format(time.RFC3339),
			"endpoints": map[string]string{
				"prometheus": "http://localhost:9090/metrics",
				"health":     "/health",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metricsHealth)
	})

	// Set CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowOriginFunc: func(origin string) bool {
			if origin == "" {
				return true
			}
			for _, allowedOrigin := range []string{"http://localhost:3000", "http://127.0.0.1:3000"} {
				if origin == allowedOrigin {
					return true
				}
			}
			return false
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	// Add middleware for HTTP requests
	handlerWithMetrics := prometheusMiddleware(mux)

	// Wrap the router with CORS middleware
	handler := c.Handler(handlerWithMetrics)

	// Start server
	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Metrics available on http://localhost:9090/metrics")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
