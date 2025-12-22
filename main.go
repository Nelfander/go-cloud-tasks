package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"task-app/taskmanager"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin) // Scanner for reading user input

	var tasks []taskmanager.Task // Slice to hold tasks

	fileData, err := os.ReadFile("tasks.json") //		 Look for the JSON file
	if err == nil {
		// If the file exists, "Unmarshal" (decode) it into our tasks slice
		json.Unmarshal(fileData, &tasks)
		fmt.Println("--- Successfully loaded saved tasks ---")
	}

	go func() { // Start a web server in a separate goroutine
		http.HandleFunc("/tasks", handleTasks) // Handle /tasks endpoint

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "index.html") // Serve the HTML file
		})
		fmt.Println("Web server starting at http://localhost:8080/")
		http.ListenAndServe(":8080", nil)
	}()

	fmt.Println("Web server live at http://localhost:8080/tasks")

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

			fmt.Println("Enter task name:")

			scanner.Scan()             // This "waits" for the user to type and hit Enter
			taskName := scanner.Text() // This captures the whole line of input with spaces and stores it in taskName

			if taskName == "" { // checks if user entered an empty task name
				fmt.Println("Error: Task title cannot be empty!")
				continue
			}

			newTask := taskmanager.Task{ // Create a new Task struct
				ID:     len(tasks) + 1, // Assign ID based on current number of tasks
				Title:  taskName,
				IsDone: false, // New tasks are not done by default
			}
			tasks = append(tasks, newTask) // Add the new task to the tasks slice
			saveTasks(tasks)               // Save tasks to my JSON file
			fmt.Println("Task added successfully!")

		} else if input == "2" {
			fmt.Println("\n---Your current Tasks---")
			for _, task := range tasks { // Iterate through tasks
				status := " " //default is empty space
				if task.IsDone {
					status = "X" //if done, mark with X

				}
				fmt.Printf("%d. [%s] %s\n", task.ID, status, task.Title) //print task details
			}

		} else if input == "3" {
			fmt.Println("Current Tasks:") // Show current tasks with IDs
			for _, t := range tasks {     // Iterate through tasks
				fmt.Printf("%d. %s\n", t.ID, t.Title)
			}
			fmt.Println("Enter the task ID to mark as done:")

			scanner.Scan()
			idInput := scanner.Text() // Read task ID input

			importID, _ := strconv.Atoi(idInput) // Convert input (string) to integer

			for i := 0; i < len(tasks); i++ { // Iterate through tasks to find the one with matching ID
				if tasks[i].ID == importID { // If found
					tasks[i].IsDone = true // Mark the task as done
					saveTasks(tasks)       // Save tasks to  my JSON file
					fmt.Printf("%v task marked as done!\n", tasks[i].Title)
				}
			}

		} else if input == "4" {
			fmt.Println("Current Tasks:") // Show current tasks with IDs
			for _, t := range tasks {     // Iterate through tasks
				fmt.Printf("%d. %s\n", t.ID, t.Title)
			}
			fmt.Print("Enter the task ID to delete: ")

			scanner.Scan()
			idInput := scanner.Text() // Read task ID input

			idToDelete, err := strconv.Atoi(idInput) // Convert input (string) to integer

			if err != nil { // Handle conversion error
				fmt.Println("Error: Please enter a valid number, not text!")
				continue
			}

			found := false
			for i := 0; i < len(tasks); i++ {
				if tasks[i].ID == idToDelete {
					tasks = append(tasks[:i], tasks[i+1:]...)
					found = true
					fmt.Println("Task deleted!")
					saveTasks(tasks) // Save tasks to my JSON file
					break            // Stop looking once we found it
				}
			}

			if !found {
				fmt.Println("Error: Task ID not found!")
			}

		} else if input == "5" {
			fmt.Println("Exiting...")
			os.Exit(0) // Exit the program
		} else {
			fmt.Println("Invalid choice, please try again.")
		}

	}

}

func saveTasks(tasks []taskmanager.Task) {
	//  Convert the slice to "Pretty" JSON
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	//  Write the data to a file named "tasks.json"
	// 0644 is the standard file permission (read/write for owner)
	err = os.WriteFile("tasks.json", data, 0644)
	if err != nil {
		fmt.Println("Error saving file:", err)
	}
}

func handleTasks(w http.ResponseWriter, r *http.Request) { // HTTP handler to serve tasks as JSON
	// Set the header so the browser knows JSON is coming
	w.Header().Set("Content-Type", "application/json")

	// Read the file and send it straight to the web browser
	fileData, _ := os.ReadFile("tasks.json")
	w.Write(fileData)
}
