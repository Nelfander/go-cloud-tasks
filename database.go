package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"task-app/taskmanager"

	_ "github.com/jackc/pgx/v5/stdlib" // Import pgx driver
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func initDB() { // Initialize the database connection
	// Load .env for local development
	_ = godotenv.Load()

	connStr := os.Getenv("DB_URL") // Get the connection string from environment variable
	var err error
	db, err = sql.Open("pgx", connStr) // Use pgx driver
	if err != nil {
		panic(err) // Handle connection error
	}
	// Ping the database to ensure the connection string is valid
	err = db.Ping()
	if err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}
	// Create BOTH tables
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username TEXT UNIQUE NOT NULL,
        password_hash TEXT NOT NULL
    );
    CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        is_done BOOLEAN DEFAULT FALSE,
        user_id INTEGER REFERENCES users(id)
    );`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println("Cloud Database connected and table ready!")
}

func RegisterUser(username, password string) error { // 	Function to register a new user
	// Hash the password (cost 10 is the industry standard)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insert into the "users" table
	// I use $1 and $2 to prevent SQL Injection
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2)`
	_, err = db.Exec(query, username, string(hashedPassword))

	return err
}

func GetUserByUsername(username string) (taskmanager.User, error) { // Function to get a user by username
	var u taskmanager.User
	query := `SELECT id, username, password_hash FROM users WHERE username = $1`

	//  QueryRow because we only expect ONE user
	err := db.QueryRow(query, username).Scan(&u.ID, &u.Username, &u.PasswordHash)
	return u, err
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

func GetTasksByUserID(userID int) ([]taskmanager.Task, error) { // Function to get tasks for a specific user
	// We filter by user_id
	rows, err := db.Query("SELECT id, title, is_done FROM tasks WHERE user_id = $1", userID)
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

func ToggleTaskSafe(taskID int, userID int) error { // Function to toggle task status with user ownership check
	// I use AND user_id = $2 to ensure ownership
	result, err := db.Exec("UPDATE tasks SET is_done = NOT is_done WHERE id = $1 AND user_id = $2", taskID, userID)
	if err != nil {
		return err
	}

	// Check if any row was actually changed
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("task not found or unauthorized")
	}

	return nil // Success
}
