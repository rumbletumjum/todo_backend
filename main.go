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

func (server *TodoServer) getAll(w http.ResponseWriter, req *http.Request) {
	log.Printf("getAll: %v", req)
	todos, err := server.store.GetAll()
	err = json.NewEncoder(w).Encode(todos)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (server *TodoServer) save(w http.ResponseWriter, req *http.Request) {
	log.Printf("save: %v", req)
	var todo todo.Todo
	err := json.NewDecoder(req.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = server.store.Save(&todo)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (server *TodoServer) handle(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		server.getAll(w, req)
	case "POST":
		server.save(w, req)
	}
}

func main() {
	mux := http.NewServeMux()
	server := NewTodoServer()

	mux.HandleFunc("/todo", server.handle)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
