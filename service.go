package main

import (
	"errors"
	"fmt"
)

var ErrQueueNotExists = errors.New("queue not exists")
var ErrNotFound = ErrEmpty

// getQueue отдает элемент из очереди
// где URL path - имя очереди
func getQueue(manager *Manager, queueName string) (string, error) {
	if !manager.Exists(queueName) {
		return "", ErrQueueNotExists
	}

	element, err := manager.Get(queueName).Get()
	if err != nil {
		return "", fmt.Errorf("getting element: %w", err)
	}

	return element, nil
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
