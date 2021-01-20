package main

import (
	"encoding/json"
	"log"
	"net/http"
	"rkb/todo_backend/internal/store"
)

type TodoServer struct {
	store store.TodoStore
}

func NewTodoServer() *TodoServer {
	store := store.NewInMemoryTodoStore()
	return &TodoServer{store: store}
}

func (server *TodoServer) getAllTodos(w http.ResponseWriter, req *http.Request) {
	log.Printf("getAll: %v", req)

	todos := server.store.GetAll()
	err := json.NewEncoder(w).Encode(todos)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (server *TodoServer) newTodo(w http.ResponseWriter, req *http.Request) {
	log.Printf("save: %v", req)

	var todo store.Todo
	err := json.NewDecoder(req.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = server.store.NewTodo(&todo)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (server *TodoServer) handle(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		server.getAllTodos(w, req)
	case "POST":
		server.newTodo(w, req)
	}
}

func main() {
	mux := http.NewServeMux()
	server := NewTodoServer()

	mux.HandleFunc("/todo", server.handle)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
