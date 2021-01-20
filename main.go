package main

import (
	"encoding/json"
	"log"
	"net/http"
	"rkb/todo_backend/internal/store"
	"rkb/todo_backend/internal/todo"
)

type TodoServer struct {
	store store.TodoStore
}

func NewTodoServer() *TodoServer {
	store := store.NewInMemoryTodoStore()
	return &TodoServer{store: store}
}

func (srv *TodoServer) getAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("getAll: %v", r)
	todos, err := srv.store.GetAll()
	err = json.NewEncoder(w).Encode(todos)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (srv *TodoServer) save(w http.ResponseWriter, r *http.Request) {
	log.Printf("save: %v", r)
	var todo todo.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = srv.store.Save(&todo)
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
	srv := NewTodoServer()

	http.HandleFunc("/todo", srv.handle)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
