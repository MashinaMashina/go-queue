package main

import (
	"sync"
)

// Manager - менеджер очередей
type Manager struct {
	mu     sync.Mutex
	queues map[string]*Queue
}

func NewManager() *Manager {
	return &Manager{
		queues: make(map[string]*Queue),
	}
}

// Get возвращает очередь по имени, если её еще нет - добавляет
func (m *Manager) Get(name string) *Queue {
	m.mu.Lock()
	defer m.mu.Unlock()

	if q, exists := m.queues[name]; exists {
		return q
	}

	m.queues[name] = NewQueue()

	return m.queues[name]
}
