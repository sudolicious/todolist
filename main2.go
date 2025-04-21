package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"net/http"
	"strconv"

        "github.com/golang-migrate/migrate/v4"
        "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

)

type Task struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Done bool `json:"done"`
	CreatedAt *time.Time `json:"created_at"`
}

func (t *Task) FormattedCreatedAt() string {
	if t.CreatedAt == nil ||
	t.CreatedAt.IsZero() {
		return "No data"
	}
	return t.CreatedAt.Format("2006-01-02 15:04")
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		done BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP)
	`)
	return err
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

func main() {
	//Connection to PostgreSql
	connStr := "user=olgadb dbname=dbgo password='Cvetaria2015' sslmode=disable"
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

	//Create table
	if err = createTable(db);
	err != nil {
		log.Fatal("Mistake of creating the table", err)
	}

	//Set migration
	driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
		log.Fatal("Mistake of driver:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", "postgres", driver)
	if err != nil {
		log.Fatal("Mistake of migration:", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Mistake of migration:", err)
	}

	// Set HTTP Router
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := getAllTasks(db)
		if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
		}
		// JSON output usinq encoding
		w.Header().Set("Content-Type","application/json")

		json.NewEncoder(w).Encode(tasks)
		})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	tasks, err := getAllTasks(db)
		if err != nil {
    		log.Fatal(err)
	}

	fmt.Println("List of tasks:")
		for _, task := range tasks {
    	fmt.Printf("%d: %s (done: %t, created: %s)\n", task.ID, task.Title, task.Done, task.FormattedCreatedAt())
}
})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/done", func(w http.ResponseWriter, r *http.Request) {
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

	if err = completeTask(db, id);
	err != nil {
		http.Error(w, err.Error(),http.StatusInternalServerError)
	return
}
	w.WriteHeader(http.StatusOK)
})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
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
	if err = deleteTask(db,id);
	err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	return
	}
	w.WriteHeader(http.StatusOK)
})

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
