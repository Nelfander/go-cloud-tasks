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

func handleCreateTask(w http.ResponseWriter, r *http.Request) { // Handler to create a new task
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Get the User ID from the context (just like in handleTasks)
	userID := r.Context().Value("userID").(int)

	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// 2. Pass the userID into your Save function
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

func handleLogin(w http.ResponseWriter, r *http.Request) { // New handler for user login
	log.Println("--- Login Attempt Started ---") // Step 0 for debugging
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
	log.Printf("Step 1: Received data for user: %s\n", username) // Step 1 for debugging
	fmt.Printf("Login attempt for user: %s\n", username)

	//  Look up user in DB
	user, err := GetUserByUsername(username)
	if err != nil {
		log.Printf("Step 2 Error: User not found or DB error: %v\n", err) // Step 2 for debugging
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	//  Compare Password with Hash from the DB
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Println("Step 3 Error: Password does not match") //		 Step 3 for debugging
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
	log.Println("Step 3: Password verified") // step 4 for debugging

	// Sign the token with the secret key
	tokenString, err := token.SignedString(getJWTSecret()) // Use the function to get secret
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	log.Println("Step 4: JWT created successfully") // Step 5 for debugging
	//  Send it back to the user
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
	log.Println("Step 5: JSON response sent to browser") // Step 6 for debugging

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
