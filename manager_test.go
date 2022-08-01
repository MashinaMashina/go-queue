package main

import (
	"testing"
)

func TestManager(t *testing.T) {
	queueName := "test-queue"
	element := "test-element"
	manager := NewManager()

	if manager.Exists(queueName) {
		t.Fatalf("queue is exists on empty manager: %s", queueName)
	}

	queue := manager.Get(queueName)
	_, err := queue.Get()

	if err == nil {
		t.Fatalf("exists elements on empty queue")
	}

	queue.Add(element)

	value, err := queue.Get()

	if err != nil {
		t.Fatalf("error on getting elements from not empty queue: %s", err.Error())
	}

	if value != element {
		t.Fatalf("put '%s', but got '%s' from queue", element, value)
	}
}
