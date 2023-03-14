package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Command line flags
var (
	directory = flag.String("directory", "", "Directory of markdown files")
)

// Enumeration of the different types of todos, used in the getTodos function
// to search for either done or uncompleted todos.
type TaskStatus int

const (
	DONE TaskStatus = iota
	NOT_DONE
)

func main() {

	flag.Parse()

	setDoneTodoInOriginalFile(*directory+"output.md", *directory)
	aggregateTodos(*directory)
}

// Aggregate all the todos from all the files in a directory.
// The code is pretty straightforward. The main function calls aggregateTodos which calls
// getMarkdownFiles to get all the markdown files in a directory and its subdirectories.
// Then it iterates through each file and calls getTodos to get all the todos from that file.
// Finally, it calls writeTodos to write the todos from all files to a file named output.md.
func aggregateTodos(directory string) {

	files, err := getMarkdownFiles(directory)
	var output string = directory + "output.md"

	// Remove the output file if it already exists, otherwise the new todos will be appended to the
	// old ones.
	os.Remove(output)

	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(files); i++ {
		todos, _ := getTodos(files[i], NOT_DONE)
		writeTodos(output, todos)
	}
}

// Get all the todos from a file and return them as a slice.
// Based on the taskStatus parameter, it will return either done or uncompleted todos.
func getTodos(filename string, taskStatus TaskStatus) ([]string, error) {

	fmt.Println("Opening file:", filename)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()

	// Create a scanner that reads from the file
	scanner := bufio.NewScanner(file)

	// Store the todos in a slice
	var todos []string

	// Iterate through each line of the file to search for done or uncompleted todos
	var prefix string
	if taskStatus == DONE {
		prefix = "- [x]"
	} else if taskStatus == NOT_DONE {
		prefix = "- [ ]"
	}
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line) + "\n"

		if strings.HasPrefix(trimmedLine, prefix) {
			todos = append(todos, trimmedLine)
		}
	}

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

// Write a slice of todos to a file.
// Example filename: "/tmp/output.md"
func writeTodos(filename string, todos []string) error {

	fmt.Println("Writing to file:", filename)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	for i := 0; i < len(todos); i++ {
		_, err = file.WriteString(todos[i])
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return err
		}
	}
	return nil
}

// Get all markdown files in a directory and its subdirectories
func getMarkdownFiles(directory string) ([]string, error) {
	var files []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking the path:", err)
		return nil, err
	}
	return files, nil
}

// Mark a completed todo back in the original file
func setDoneTodoInOriginalFile(outputFile string, directory string) {

	// 1. get all todos from the output file
	todos, _ := getTodos(outputFile, DONE)
	// 2. For all line starting with "- [x]"
	for i := 0; i < len(todos); i++ {
		line := todos[i]
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "- [x]") {
			// 3. replace "- [x]" with "- [x]" in the original file
			file := findTodoInFile(trimmedLine, directory)
			replaceTodoInFile(trimmedLine, file)
		}
	}
}

// Replace a todo in the original file with the completed todo from the output file.
func replaceTodoInFile(todo string, filename string) {
	fmt.Println("Replacing todo in file:", filename)
	// Define the file path and the search and replacement strings
	task := strings.Replace(todo, "[x]", "[ ]", 1)

	// Read the entire file into a byte slice
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the byte slice to a string and split it into lines
	lines := strings.Split(string(content), "\n")

	// Find the index of the line that contains the search string
	var index int = -1
	for i, line := range lines {
		if strings.Contains(line, task) {
			index = i
			break
		}
	}

	// If the search string was found, replace the line with the replacement string
	if index != -1 {
		lines[index] = strings.Replace(lines[index], task, todo, 1)
	} else {
		fmt.Printf("Could not find \"%s\" in file \"%s\"\n", task, filename)
		return
	}

	// Join the lines into a single string and write it back to the file
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filename, []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Replaced \"%s\" with \"%s\" in file \"%s\"\n", task, todo, filename)
}

// Find in which file a todo is located.
func findTodoInFile(todo string, directory string) string {
	fmt.Println("Searching original file for:", todo)
	files, _ := getMarkdownFiles(directory)

	// Mark the todo as not done, so that we can find it in the original file
	todo = strings.ReplaceAll(todo, "- [x]", "- [ ]")

	// Iterate through all files in the directory
	for i := 0; i < len(files); i++ {
		// If file is the output file, skip it
		if strings.HasSuffix(files[i], "output.md") {
			continue
		}
		file, err := os.Open(files[i])
		if err != nil {
			fmt.Println("Error opening file:", err)
			return ""
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if line == todo {
				fmt.Println("Task found in file", files[i])
				return files[i]
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error scanning file:", err)
			return ""
		}
		//print as verbose
		//fmt.Println("Line not found in file", files[i])
	}
	return ""
}
