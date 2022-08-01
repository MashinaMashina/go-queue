package main

import (
	"context"
	"sync"
	"time"
)

type WaitResp struct {
	Value string
	Err   error
}

type Waiter struct {
	Result chan WaitResp
	Ctx    context.Context
	Cancel context.CancelFunc
}

type Queue struct {
	mu       sync.Mutex
	elements []string
	waiters  []*Waiter
}

// NewQueue создает очередь
func NewQueue() *Queue {
	return &Queue{
		mu:       sync.Mutex{},
		elements: make([]string, 0),
	}
}

// AddWaiter добавляет ожидающего клиента в пул
func (q *Queue) AddWaiter(duration time.Duration, client chan<- string) {
	go func() {
		q.mu.Lock()

		if len(q.elements) > 0 {
			var result string
			// берем последний элемент и вырезаем его
			result, q.elements = q.elements[len(q.elements)-1], q.elements[:len(q.elements)-1]

			client <- result
			q.mu.Unlock()
			return
		}

		// если таймаут не установлен - сразу отдаем ответ
		if duration == 0 {
			client <- ""
			q.mu.Unlock()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), duration)

		w := &Waiter{
			Result: make(chan WaitResp),
			Cancel: cancel,
			Ctx:    ctx,
		}

		q.waiters = append(q.waiters, w)

		q.mu.Unlock()

		select {
		case <-w.Ctx.Done():
			// ожидающему клиенту отправляем ошибку
			client <- ""

		case value, ok := <-w.Result:
			if ok {
				client <- value.Value
			}
		}

		q.mu.Lock()
		for i, _ := range q.waiters {
			// отработанный ожидатель удаляем
			if q.waiters[i] == w {
				q.waiters[i] = nil
				break
			}
		}

		// ищем последовательность отработанных ожидателей
		canRemove := -1
		for i, _ := range q.waiters {
			if q.waiters[i] != nil {
				break
			}
			canRemove = i
		}

		// если последовательность найдена - удаляем
		if canRemove >= 0 {
			q.waiters = q.waiters[canRemove+1:]
		}
		q.mu.Unlock()
	}()
}

// Add добавляет элемент в очередь
func (q *Queue) Add(e string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.waiters) > 0 {
		var waiter *Waiter
		waiter, q.waiters = q.waiters[0], q.waiters[1:]

		waiter.Result <- WaitResp{Value: e}
		waiter.Cancel()

		return
	}

	q.elements = append(q.elements, e)
}
