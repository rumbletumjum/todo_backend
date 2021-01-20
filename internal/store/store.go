package store

import (
	"fmt"
	"log"
	"sync"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TodoStore interface {
	GetAll() []Todo
	NewTodo(todo *Todo) error
}

type InMemoryTodoStore struct {
	sync.Mutex
	nextId int
	Todos  map[int]Todo
}

func (s *InMemoryTodoStore) GetAll() []Todo {
	log.Printf("store: GetAll")
	s.Lock()
	defer s.Unlock()

	allTodos := make([]Todo, 0, len(s.Todos))
	for _, todo := range s.Todos {
		allTodos = append(allTodos, todo)
	}
	return allTodos

}

func (s *InMemoryTodoStore) NewTodo(todo *Todo) error {
	log.Printf("store: NewTodo %v", todo)
	s.Lock()
	defer s.Unlock()

	if todo.ID == 0 {
		todo.ID = s.nextId
		s.nextId++
		s.Todos[todo.ID] = *todo
		return nil
	}
	return fmt.Errorf("unable to Save")
}

func NewInMemoryTodoStore() *InMemoryTodoStore {
	s := new(InMemoryTodoStore)
	s.Lock()
	defer s.Unlock()
	s.Todos = make(map[int]Todo)
	s.nextId = 1
	return s
}
