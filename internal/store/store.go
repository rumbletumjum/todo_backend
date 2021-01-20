package store

import (
	"fmt"
	"log"
	"rkb/todo_backend/internal/todo"
	"sync"
)

type TodoStore interface {
	GetAll() ([]*todo.Todo, error)
	Save(todo *todo.Todo) error
}

type InMemoryTodoStore struct {
	sync.Mutex
	nextId int
	Todos  []*todo.Todo
}

func (s *InMemoryTodoStore) GetAll() ([]*todo.Todo, error) {
	log.Printf("store: GetAll")
	return s.Todos, nil
}

func (s *InMemoryTodoStore) Save(todo *todo.Todo) error {
	log.Printf("store: Save %v", todo)
	s.Lock()
	defer s.Unlock()
	if todo.ID == 0 {
		todo.ID = s.nextId
		s.nextId++
		s.Todos = append(s.Todos, todo)
		return nil
	}
	return fmt.Errorf("unable to Save")
}

func NewInMemoryTodoStore() *InMemoryTodoStore {
	s := new(InMemoryTodoStore)
	s.Lock()
	defer s.Unlock()
	s.Todos = make([]*todo.Todo, 0)
	s.nextId = 1
	return s
}
