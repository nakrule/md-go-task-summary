package main

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func TestGetTodos(t *testing.T) {

	// Test case 1: Empty file
	filename := "/tmp/empty-file.md"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	expected1 := []string{}
	if output1 := GetTodos(filename, DONE); !reflect.DeepEqual(output1, expected1) {
		t.Errorf("Test case 1 failed: expected %v but got %v", expected1, output1)
	}

	// Test case 2: File with no todos
	noTodosFile := "This is just some text with no todos"
	expected2 := []string{}
	if output2 := GetTodos(noTodosFile); !reflect.DeepEqual(output2, expected2) {
		t.Errorf("Test case 2 failed: expected %v but got %v", expected2, output2)
	}

	// Test case 3: File with one todo
	oneTodoFile := "# Todo\n- [ ] Do something\nThis is some more text"
	expected3 := []string{"Do something"}
	if output3 := GetTodos(oneTodoFile); !reflect.DeepEqual(output3, expected3) {
		t.Errorf("Test case 3 failed: expected %v but got %v", expected3, output3)
	}

	// Test case 4: File with multiple todos
	multiTodoFile := "# Todos\n- [ ] Do something\n- [x] Do something else\n- [ ] Do something more\n"
	expected4 := []string{"Do something", "Do something more"}
	if output4 := GetTodos(multiTodoFile); !reflect.DeepEqual(output4, expected4) {
		t.Errorf("Test case 4 failed: expected %v but got %v", expected4, output4)
	}
}
