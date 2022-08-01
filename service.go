package main

import (
	"errors"
	"fmt"
	"time"
)

var ErrNotFound = errors.New("queue is empty")

// getQueue отдает элемент из очереди
// где URL path - имя очереди
func getQueue(manager *Manager, queueName string, d time.Duration) (string, error) {
	ch := make(chan string)
	manager.Get(queueName).AddWaiter(d, ch)
	resp := <-ch

	if resp == "" {
		return "", fmt.Errorf("getting element: %w", ErrNotFound)
	}

	return resp, nil
}

// putQueue добавляет данные в очередь
// где URL path - имя очереди, параметр v - добавляемые данные
func putQueue(manager *Manager, queueName, element string) error {
	if queueName == "" || element == "" {
		return fmt.Errorf("empty queue or element")
	}

	manager.Get(queueName).Add(element)
	return nil
}
