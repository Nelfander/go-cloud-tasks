package main

import (
	"database/sql"
	"fmt"
	"os"
	"task-app/taskmanager"

	_ "github.com/jackc/pgx/v5/stdlib" // Import pgx driver
	"github.com/joho/godotenv"
)

var db *sql.DB

func initDB() { // Initialize the database connection
	// This looks for the .env file and loads it into the system
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	connStr := os.Getenv("DB_URL") // Get the connection string from environment variable

	db, err = sql.Open("pgx", connStr) // Use pgx driver
	if err != nil {
		panic(err) // Handle connection error
	}

	// This part stays the same—it creates the table in the cloud!
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		is_done BOOLEAN DEFAULT FALSE
	);`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println("Cloud Database connected and table ready!")
}

// GetAllTasks fetches everything from the cloud

// function to not reuse the same code for showing all tasks
func GetAllTasks() ([]taskmanager.Task, error) {
	rows, err := db.Query("SELECT id, title, is_done FROM tasks ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []taskmanager.Task
	for rows.Next() {
		var t taskmanager.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.IsDone); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// AddTask saves a new task
func AddTask(title string) error {
	_, err := db.Exec("INSERT INTO tasks (title, is_done) VALUES ($1, $2)", title, false)
	return err
}

// UpdateTaskStatus toggles completion
func UpdateTaskStatus(id int, status bool) error {
	_, err := db.Exec("UPDATE tasks SET is_done = $1 WHERE id = $2", status, id)
	return err
}

// DeleteTask removes a row
func DeleteTask(id int) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
	return err
}
