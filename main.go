package main

import (
	"fmt"
	"net/http"
)

var idMap = make(map[int]int) // Map for task IDs (To show the user correct IDs after deletions )

func main() {

	// Add this line!
	//http.Handle("/", http.FileServer(http.Dir("./")))

	initDB() // Initialize the database connection

	go func() { // Start a web server in a separate goroutine
		http.HandleFunc("/tasks", handleTasks)       // Handle /tasks endpoint
		http.HandleFunc("/update", handleUpdateTask) // Handle /update endpoint
		http.HandleFunc("/create", handleCreateTask) // Handle /create endpoint
		http.HandleFunc("/delete", handleDeleteTask) // Handle /delete endpoint

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "index.html") // Serve the HTML file
		})
		fmt.Println("Web server starting at http://16.171.16.175:8080/")
		http.ListenAndServe(":8080", nil)
	}()

	fmt.Println("🚀 Application is running. Visit http://localhost:8081")
	select {} // Block forever
}

func showTasksAndPopulateMap(idMap map[int]int) { // Helper function to show tasks and fill the ID map
	tasks, err := GetAllTasks()
	if err != nil {
		fmt.Println("❌ Error fetching tasks:", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("\n📝 Your list is currently empty.")
		return
	}

	fmt.Printf("\n%-5s %-20s %-10s\n", "#", "Title", "Status")
	fmt.Println("---------------------------------------")

	for i, t := range tasks {
		displayNum := i + 1
		idMap[displayNum] = t.ID // This maps 1, 2, 3 to the real DB ID

		status := "❌"
		if t.IsDone {
			status = "✅"
		}
		fmt.Printf("%-5d %-20s %-10s\n", displayNum, t.Title, status)
	}
}
