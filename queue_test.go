package main

import (
	"testing"
)

// TestManager проверяет получение записей из пустой очереди, запись в очередь,
// получение элементов из не пустой очереди
func TestQueue(t *testing.T) {
	element := "test-element"
	queue := NewQueue()

	ch := make(chan string)
	queue.AddWaiter(0, ch)
	res := <-ch

	if res != "" {
		t.Fatalf("exists elements on empty queue")
	}

	queue.Add(element)

	ch = make(chan string)
	queue.AddWaiter(0, ch)
	res = <-ch

	if res != element {
		t.Fatalf("put '%s', but got '%s' from queue", element, res)
	}
}
