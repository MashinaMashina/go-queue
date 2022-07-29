package main

import (
	"errors"
	"sync"
)

var ErrEmpty = errors.New("queue is empty")

type Queue struct {
	mu       sync.RWMutex
	elements []string
}

// NewQueue создает очередь
func NewQueue() *Queue {
	return &Queue{
		mu:       sync.RWMutex{},
		elements: make([]string, 0),
	}
}

// Get отдает самый старый элемент из очереди
func (q *Queue) Get() (string, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if len(q.elements) == 0 {
		return "", ErrEmpty
	}

	var result string
	result, q.elements = q.elements[0], q.elements[1:]

	return result, nil
}

// Add добавляет элемент в очередь
func (q *Queue) Add(e string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.elements = append(q.elements, e)
}
