package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type InMemoryTodoService struct {
	sync.Mutex
	nextId int
	Todos  []*Todo
}

func NewTodoService() *InMemoryTodoService {
	s := new(InMemoryTodoService)
	s.Lock()
	defer s.Unlock()
	s.Todos = make([]*Todo, 0)
	s.nextId = 1
	return s
}

func (s *InMemoryTodoService) getAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("getAll: %v", r)
	err := json.NewEncoder(w).Encode(s.Todos)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (s *InMemoryTodoService) save(w http.ResponseWriter, r *http.Request) {
	log.Printf("save: %v", r)
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if todo.ID == 0 {
		s.Lock()
		defer s.Unlock()
		todo.ID = s.nextId
		s.nextId++

		s.Todos = append(s.Todos, &todo)
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (s *InMemoryTodoService) handle(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        s.getAll(w, r)
    case "POST":
        s.save(w, r)
    }
}

func main() {
	svc := NewTodoService()

	http.HandleFunc("/todo", svc.handle)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
