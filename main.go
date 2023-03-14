package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	directory = flag.String("directory", "", "Directory of markdown files")
)

func main() {

	flag.Parse()

	files, err := getMarkdownFiles(*directory)
	var output string = *directory + "output.md"

	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(files); i++ {
		todos, _ := getTodos(files[i])
		writeTodos(output, todos)
	}
}

// Get all the todos from a file
func getTodos(filename string) ([]string, error) {
	// Open the file for reading
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

	// Iterate through each line of the file
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line) + "\n"

		// Check if the line starts with "- [ ]"
		if strings.HasPrefix(trimmedLine, "- [ ]") {
			todos = append(todos, trimmedLine)
		}
	}

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

// Write the todos to a file
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
func getMarkdownFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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
