package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func handleTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Use  function from database.go
	tasks, err := GetAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// Send the results back as JSON
	json.NewEncoder(w).Encode(tasks)

	// Get the ID from the URL (e.g., /update?id=10)
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	err = UpdateTaskStatus(id, true) // Using  clean function!
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the ID from the URL (e.g., /update?id=10)
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr) // Convert to integer

	err := UpdateTaskStatus(id, true) // Using your clean function!
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	// We expect a simple query parameter ?title=Work
	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "Title is required", 400)
		return
	}

	err := AddTask(title) // clean function from database.go
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	err := DeleteTask(id) //  clean function from database.go
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleRegister(w http.ResponseWriter, r *http.Request) { // New handler for user registration
	if r.Method == http.MethodGet {
		// Show the registration form
		http.ServeFile(w, r, "register.html")
		return
	}

	if r.Method == http.MethodPost {
		// Grab values from the HTML form
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Use the logic from database.go to save to Neon
		err := RegisterUser(username, password)
		if err != nil {
			// If the username is already in the DB, this will trigger
			http.Error(w, "Registration failed. Username might be taken.", http.StatusBadRequest)
			return
		}

		// Success message
		fmt.Fprintf(w, "Success! User %s created. Happy Christmas Eve!", username)
	}
}
