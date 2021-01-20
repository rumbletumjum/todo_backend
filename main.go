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

type TodoStore interface {
	getAll() ([]*Todo, error)
	save(todo *Todo) error
}

type InMemoryTodoStore struct {
	sync.Mutex
	nextId int
	Todos  []*Todo
}

func (s *InMemoryTodoStore) getAll() ([]*Todo, error) {
	log.Printf("store: getAll")
	return s.Todos, nil
}

func (s *InMemoryTodoStore) save(todo *Todo) error {
	log.Printf("store: save %v", todo)
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

func NewTodoStore() *InMemoryTodoStore {
	s := new(InMemoryTodoStore)
	s.Lock()
	defer s.Unlock()
	s.Todos = make([]*Todo, 0)
	s.nextId = 1
	return s
}

type TodoServer struct {
	store TodoStore
}

func (srv *TodoServer) getAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("getAll: %v", r)
	todos, err := srv.store.getAll()
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
	err = srv.store.save(&todo)
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
	srv := TodoServer{store: NewTodoStore()}

	http.HandleFunc("/todo", srv.handle)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
