package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

type Task struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Done      bool       `json:"done"`
	CreatedAt *time.Time `json:"created_at"`
}

func addTask(db *sql.DB, title string) (*Task, error) {
	var task Task
	err := db.QueryRow(
		"INSERT INTO tasks (title) VALUES ($1) RETURNING id, title, done, created_at", title,
	).Scan(&task.ID, &task.Title, &task.Done, &task.CreatedAt)
	return &task, err
}

func getAllTasks(db *sql.DB) ([]Task, error) {
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

func completeTask(db *sql.DB, id int) error {
	_, err := db.Exec("UPDATE tasks SET done = TRUE WHERE id = $1", id)
	return err
}

func deleteTask(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
	return err
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

	//Connection to PostgreSql
	connStr := "user=olgadb dbname=dbgo password='****' sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to PostgreSql!")

	//Migration func
	if err := runMigrations(db); err != nil {
		log.Fatal("Migration error:", err)
	}
	fmt.Println("Migration applied successfully")

	//Create Router
        mux := http.NewServeMux()

	// Set HTTP Router
	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := getAllTasks(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// JSON output usinq encoding
		w.Header().Set("Content-Type",	"application/json",)

		// cashing
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

		task, err := addTask(db, title)
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
		// get id from URL and convert string
		idStr := r.FormValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		if err = completeTask(db, id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
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
		if err = deleteTask(db, id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	//Set CORS
	c := cors.New(cors.Options{
    		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
    		AllowOriginFunc: func(origin string) bool {

		// Allow requests with no Origin (for curl)
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
    	Debug:           true,
	})

	// Wrap the router with CORS middleware
    	handler := c.Handler(mux)

	fs := http.FileServer(http.Dir("./frontend/build"))
	mux.Handle("/", http.StripPrefix("/", fs))

	// Start server
	fmt.Println("Server running on http://localhost:8080")
    	log.Fatal(http.ListenAndServe(":8080", handler))
}
