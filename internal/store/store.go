package store

import (
	"log"
	"sync"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TodoStore interface {
	GetAllTodos() []Todo
	NewTodo(title string) int
}

type InMemoryTodoStore struct {
	sync.Mutex
	nextId int
	Todos  map[int]Todo
}

func (s *InMemoryTodoStore) GetAllTodos() []Todo {
	log.Printf("store: GetAll")
	s.Lock()
	defer s.Unlock()

	allTodos := make([]Todo, 0, len(s.Todos))
	for _, todo := range s.Todos {
		allTodos = append(allTodos, todo)
	}
	return allTodos
}

func (s *InMemoryTodoStore) NewTodo(title string) int {
	s.Lock()
	defer s.Unlock()

	todo := Todo{
		ID:        s.nextId,
		Title:     title,
		Completed: false,
	}

	s.Todos[s.nextId] = todo
	s.nextId++
	return todo.ID
}

func NewInMemoryTodoStore() *InMemoryTodoStore {
	s := new(InMemoryTodoStore)
	s.Lock()
	defer s.Unlock()
	s.Todos = make(map[int]Todo)
	s.nextId = 1
	return s
}
