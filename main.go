package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func main() {

	// Add this line!
	//http.Handle("/", http.FileServer(http.Dir("./")))

	initDB()                              // Initialize the database connection
	scanner := bufio.NewScanner(os.Stdin) // Scanner for reading user input

	var idMap = make(map[int]int) // Map for task IDs (To show the user correct IDs after deletions )

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

	fmt.Println("Web server live at http://16.171.16.175:8080/")

	for { // Infinite loop for the menu
		fmt.Println("\n---GO-Task Manager---")
		fmt.Println("1. Add Task")
		fmt.Println("2. List of Tasks")
		fmt.Println("3. Mark Task as Done")
		fmt.Println("4. Delete Task")
		fmt.Println("5. Exit")
		fmt.Print("Choose an option: ")

		scanner.Scan()          // This "waits" for the user to type and hit Enter
		input := scanner.Text() // Read user input

		if input == "1" {

			fmt.Print("Enter task: ")
			scanner.Scan()
			AddTask(scanner.Text()) // Just one clean line!
			fmt.Println("Done.")

		} else if input == "2" {
			showTasksAndPopulateMap(idMap) //	 Helper function to show tasks

		} else if input == "3" {
			showTasksAndPopulateMap(idMap)
			fmt.Print("\nEnter the # to mark as DONE: ")
			var choice int
			fmt.Scan(&choice)

			if realID, exists := idMap[choice]; exists {
				UpdateTaskStatus(realID, true)
				fmt.Println("✅ Task updated!")
			} else {
				fmt.Println("⚠️ Invalid selection.")
			}

		} else if input == "4" {
			showTasksAndPopulateMap(idMap)
			fmt.Print("\nEnter the # to DELETE: ")
			var choice int
			fmt.Scan(&choice)

			if realID, exists := idMap[choice]; exists {
				DeleteTask(realID)
				fmt.Println("🗑️ Task removed!")
			} else {
				fmt.Println("⚠️ Invalid selection.")
			}

		} else if input == "5" {
			fmt.Println("Closing database connection...")
			db.Close() // Good practice to close the DB before the app shuts down
			fmt.Println("Goodbye! 👋")
			return // Exits the main function (and the program)

		} else {
			fmt.Println("Invalid option, please try again.")
		}

	}

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
