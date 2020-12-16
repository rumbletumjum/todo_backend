package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TodoService interface {
	GetAll() ([]*Todo, error)
	Save(todo *Todo) error
}

type InMemoryTodoService struct {
	mu     sync.Mutex
	nextId int
	Todos  []*Todo
}

func NewTodoService() *InMemoryTodoService {
	s := new(InMemoryTodoService)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Todos = make([]*Todo, 0)
	s.nextId = 1
	return s
}

func (s *InMemoryTodoService) GetAll() ([]*Todo, error) {
	return s.Todos, nil
}

func (s *InMemoryTodoService) Save(todo *Todo) error {
	if todo.ID == 0 {
		s.mu.Lock()
		defer s.mu.Unlock()
		todo.ID = s.nextId
		s.nextId++

		s.Todos = append(s.Todos, todo)
		return nil
	}

	return fmt.Errorf("Unable to save")
}

var TodoSvc TodoService

func getHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := TodoSvc.GetAll()
	err = json.NewEncoder(w).Encode(todos)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = TodoSvc.Save(&todo)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func main() {
	TodoSvc = NewTodoService()

	testTodo := &Todo{
		ID:        0,
		Title:     "Finish ToDo app",
		Completed: false,
	}
	TodoSvc.Save(testTodo)
	http.HandleFunc("/todos", getHandler)
	http.HandleFunc("/todo", saveHandler)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
