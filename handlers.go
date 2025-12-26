package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func getJWTSecret() []byte { //	 Function to get JWT secret from .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: No .env file found, using system environment variables")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// WHAT IF: The secret is missing?
		log.Fatal("CRITICAL: JWT_SECRET not set in .env file")
	}
	return []byte(secret)
}
func handleTasks(w http.ResponseWriter, r *http.Request) {
	//  Method check is good! Keep it.
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//  NEW: Extract the UserID saved in Middleware context
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Could not identify user", http.StatusInternalServerError)
		return
	}

	//  UPDATED: Call a function that only gets tasks for THIS user
	tasks, err := GetTasksByUserID(userID)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//  Send the JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)

}

func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	// Get the UserID from the JWT context
	userID := r.Context().Value("userID").(int)

	//  Get the Task ID from the URL
	idStr := r.URL.Query().Get("id")
	taskID, _ := strconv.Atoi(idStr) // Convert to int

	//  Pass BOTH to the database
	err := ToggleTaskSafe(taskID, userID)
	if err != nil {
		http.Error(w, "Could not update task: "+err.Error(), http.StatusForbidden) // 403 if not allowed
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleCreateTask(w http.ResponseWriter, r *http.Request) { // Handler to create a new task
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//  Get the User ID from the context (just like in handleTasks)
	userID := r.Context().Value("userID").(int)

	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	//  Pass the userID into your Save function
	err := SaveTask(title, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func SaveTask(title string, userID int) error { // Updated Save function to include userID
	// add user_id to the INSERT statement
	query := `INSERT INTO tasks (title, is_done, user_id) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, title, false, userID)
	return err
}

func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodPost {
	//	http.Error(w, "Method not allowed", 405)
	//	return
	//}

	userID := r.Context().Value("userID").(int) // Get UserID from context
	idStr := r.URL.Query().Get("id")            // Get Task ID from URL
	taskID, _ := strconv.Atoi(idStr)            // Convert to int

	// Safe delete
	result, err := db.Exec("DELETE FROM tasks WHERE id = $1 AND user_id = $2", taskID, userID) // Ensure ownership
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rowsAffected, _ := result.RowsAffected() // Check if any row was deleted
	if rowsAffected == 0 {
		http.Error(w, "Unauthorized to delete this task", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleRegister(w http.ResponseWriter, r *http.Request) { // New handler for user registration
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "register.html")
		return
	}

	// Parse form data from the fetch request
	r.ParseMultipartForm(10 << 20)
	username := r.FormValue("username")
	password := r.FormValue("password")

	//  Validate input
	if len(username) < 3 || len(password) < 6 {
		http.Error(w, "Username (3+) or password (6+) too short", http.StatusBadRequest)
		return
	}

	//  Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	//  Insert into Database
	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, string(hashedPassword))

	if err != nil {
		// handle the "Username taken" error
		// Different drivers use different error checks, but checking for
		// string "unique_violation" or "duplicate key" is a safe bet for now.
		log.Printf("Register error: %v", err)
		http.Error(w, "That username is already taken. Try another!", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Success")
}

func handleLogin(w http.ResponseWriter, r *http.Request) { // New handler for user login

	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "login.html")
		return
	}

	// IMPORTANT: This line parses the data sent by 'new FormData(loginForm)'
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		// Fallback to standard form parsing if not multipart
		r.ParseForm()
	}

	//  Get credentials from form
	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Printf("Login attempt for user: %s\n", username)

	//  Look up user in DB
	user, err := GetUserByUsername(username)
	if err != nil {

		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	//  Compare Password with Hash from the DB
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {

		// WHAT IF: The password is wrong?
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	//  Create the Claims (the info inside the ID card)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(getJWTSecret()) // Use the function to get secret
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	//  Send it back to the user
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"token":    tokenString,
		"username": username, // update username to display on the dashboard
	})

}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { // Middleware to protect routes
		authHeader := r.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return getJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//  Get the claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Extract UserID and put it into the Request Context
		userID := int(claims["user_id"].(float64)) // JWT numbers are float64 by default
		ctx := context.WithValue(r.Context(), "userID", userID)

		//  Pass the new context to the next handler
		next(w, r.WithContext(ctx))

	}
}
