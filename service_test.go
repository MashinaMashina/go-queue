package main

import (
	"errors"
	"testing"
)

func TestGetQueue(t *testing.T) {
	queueName := "test-queue"
	element := "test-element-get"
	manager := NewManager()
	queue := manager.Get(queueName)

	queue.Add(element)

	res, err := getQueue(manager, queueName)
	if err != nil {
		t.Fatalf("got error in getting element from queue: %s", err.Error())
	}

	if res != element {
		t.Fatalf("got invalid element from queue. Expected: '%s', but real: '%s'", element, res)
	}

	_, err = queue.Get()

	if err == nil {
		t.Fatalf("error expected in getting element from empty queue")
	}

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound error. Expected: '%s', but real: '%s'", ErrNotFound.Error(), err.Error())
	}
}

func TestPutQueue(t *testing.T) {
	queueName := "test-queue"
	element := "test-element-put"
	manager := NewManager()
	queue := manager.Get(queueName)

	err := putQueue(manager, queueName, element)
	if err != nil {
		t.Fatalf("got error in putting element to queue: %s", err.Error())
	}

	res, err := queue.Get()
	if err != nil {
		t.Fatalf("got error in getting element from queue: %s", err.Error())
	}

	if res != element {
		t.Fatalf("got invalid element from queue. Expected: '%s', but real: '%s'", element, res)
	}
}
