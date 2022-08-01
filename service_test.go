package main

import (
	"testing"
	"time"
)

// TestGetQueue пытается получить один добавленный заранее элемент из очереди.
// После проверяет что очередь очистилась
func TestGetQueue(t *testing.T) {
	queueName := "test-queue"
	element := "test-element-get"
	manager := NewManager()
	queue := manager.Get(queueName)

	queue.Add(element)

	res, err := getQueue(manager, queueName, 0)
	if err != nil {
		t.Fatalf("got error in getting element from queue: %s", err.Error())
	}

	if res != element {
		t.Fatalf("got invalid element from queue. Expected: '%s', but real: '%s'", element, res)
	}

	ch := make(chan string)
	manager.Get(queueName).AddWaiter(0, ch)
	res = <-ch

	if res != "" {
		t.Fatalf("expected empty result in getting element from empty queue")
	}
}

// TestGetQueueWithTimeoutPositive добавляет ожидающего клиента и через время добавляет сообщение в очередь
// Ожидающий клиент должен получить сообщение, которое появилось в очереди после начала ожидания
func TestGetQueueWithTimeoutPositive(t *testing.T) {
	queueName := "test-queue"
	element := "test-element-get"
	manager := NewManager()
	queue := manager.Get(queueName)

	go func() {
		time.Sleep(time.Millisecond)
		queue.Add(element)
	}()

	ch := make(chan string)
	queue.AddWaiter(time.Millisecond*3, ch)
	res := <-ch

	if res != element {
		t.Fatalf("got invalid element from queue. Expected: '%s', but real: '%s'", element, res)
	}
}

// TestGetQueueWithTimeoutPositive добавляет ожидающего клиента.
// Сообщений в очередь не приходит - должен сработать таймаут.
func TestGetQueueWithTimeoutNegative(t *testing.T) {
	queueName := "test-queue"
	element := "test-element-get"
	manager := NewManager()
	queue := manager.Get(queueName)

	queue.Add(element)

	res, err := getQueue(manager, queueName, 0)
	if err != nil {
		t.Fatalf("got error in getting element from queue: %s", err.Error())
	}

	if res != element {
		t.Fatalf("got invalid element from queue. Expected: '%s', but real: '%s'", element, res)
	}

	ch := make(chan string)
	queue.AddWaiter(time.Millisecond, ch)
	res = <-ch

	if res != "" {
		t.Fatalf("expected empty result in getting element from empty queue")
	}
}

// TestPutQueue добавляет сообщение в очередь и после проверяет сообщение в очереди
func TestPutQueue(t *testing.T) {
	queueName := "test-queue"
	element := "test-element-put"
	manager := NewManager()
	queue := manager.Get(queueName)

	err := putQueue(manager, queueName, element)
	if err != nil {
		t.Fatalf("got error in putting element to queue: %s", err.Error())
	}

	ch := make(chan string)
	queue.AddWaiter(0, ch)
	res := <-ch

	if res != element {
		t.Fatalf("got invalid element from queue. Expected: '%s', but real: '%s'", element, res)
	}
}
