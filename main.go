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
	getAll() ([]*Todo, error)
	save(todo *Todo) error
}

type InMemoryTodoService struct {
	sync.Mutex
	nextId int
	Todos  []*Todo
}

func (s *InMemoryTodoService) getAll() ([]*Todo, error) {
	log.Printf("svc: getAll")
	return s.Todos, nil
}

func (s *InMemoryTodoService) save(todo *Todo) error {
	log.Printf("svc: save %v", todo)
	s.Lock()
	defer s.Unlock()
	if todo.ID == 0 {
		todo.ID = s.nextId
		s.nextId++
		s.Todos = append(s.Todos, todo)
		return nil
	}
	return fmt.Errorf("unable to save")
}

func NewTodoService() *InMemoryTodoService {
	s := new(InMemoryTodoService)
	s.Lock()
	defer s.Unlock()
	s.Todos = make([]*Todo, 0)
	s.nextId = 1
	return s
}

type TodoServer struct {
	svc TodoService
}

func (srv *TodoServer) getAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("getAll: %v", r)
	todos, err := srv.svc.getAll()
	err = json.NewEncoder(w).Encode(todos)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (srv *TodoServer) save(w http.ResponseWriter, r *http.Request) {
	log.Printf("save: %v", r)
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = srv.svc.save(&todo)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (srv *TodoServer) handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		srv.getAll(w, r)
	case "POST":
		srv.save(w, r)
	}
}

func main() {
	srv := TodoServer{svc: NewTodoService()}

	http.HandleFunc("/todo", srv.handle)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
