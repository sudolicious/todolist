package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"encoding/json"
	_ "github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Task struct {
	ID 		int		`json:"id"`
	Title		string		`json:"title"`
	Done		bool		`json:"done"`
	CreatedAt	*time.Time	`json:"created_at"` // * for null
}

//Formating time
func (t Task) FormattedCreatedAt() string {
	if t.CreatedAt == nil {
	return "No data"
	}
	return t.CreatedAt.Format("2006-01-02 15:04")
	}

func createTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		done BOOLEAN DEFAULT FALSE)
	`)
	return err
}

func addTask(db *sql.DB, title string) error {
	_, err := db.Exec("INSERT INTO tasks (title, created_at) VALUES ($1, $2)", title, time.Now())
	return err
}

func getAllTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT id, title, done, created_at FROM tasks ORDER BY id")
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
	err = createTable(db)
	if  err != nil {
		log.Fatal(err)
	}

	//Migration
	err = migrate(db)
	if err != nil {
		log.Fatal("Migration failed", err)
	}

	//Testing
	err = addTask(db, "First Task")
	if err != nil {
		log.Fatal(err)
	}

	tasks, err := getAllTasks(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("List of tasks:")
	for _, task := range tasks {
		fmt.Printf("%d: %s (done: %t)\n", task.ID, task.Title, task.Done)
}

	// Set HTTP Router
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := getAllTasks(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	type TaskResponse struct {
		ID int `json:"id"`
		Title string `json:"title"`
		Done bool `json:"done"`
		CreatedAt string `json:"created_at"`
}

	var response []TaskResponse
	for _, task := range tasks { 
		response = append(response, TaskResponse {
		ID: task.ID,
		Title: task.Title,
		Done: task.Done,
		CreatedAt: task.FormattedCreatedAt(),
		})
	}

		// JSON output
		w.Header().Set("Content-Type","application/json")
		json.NewEncoder(w).Encode(response)
})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
		}

	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
		}

	err = addTask(db, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	})

	http.HandleFunc("/done", func(w http.ResponseWriter, r *http.Request) {
	// get id from URL and convert string
	idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
}

	err = completeTask(db, id)
	if err != nil {
		http.Error(w, err.Error(),http.StatusInternalServerError)
	return
}

	w.WriteHeader(http.StatusOK)
})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = deleteTask(db,id)
		if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	return
	}

	w.WriteHeader(http.StatusOK)
})

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
