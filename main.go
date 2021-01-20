package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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
	todos := server.store.GetAllTodos()
	res, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	w.Write(res)
}

func (server *TodoServer) newTodo(w http.ResponseWriter, req *http.Request) {
	type RequestTodo struct {
		Title string `json:"title"`
	}

	type ResponseId struct {
		Id int `json:"id"`
	}

	var reqTodo RequestTodo
	err := json.NewDecoder(req.Body).Decode(&reqTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := server.store.NewTodo(reqTodo.Title)
	responseStruct := ResponseId{id}
	res, err := json.Marshal(responseStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(res)
}

func (server *TodoServer) handle(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		server.getAllTodos(w, req)
	case http.MethodPost:
		server.newTodo(w, req)
	}
}

func main() {
	mux := http.NewServeMux()
	server := NewTodoServer()

	mux.HandleFunc("/todo", server.handle)

	port, ok := os.LookupEnv("SERVEPORT")
	if !ok {
		port = "8888"
	}
	log.Println("Listening on " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
