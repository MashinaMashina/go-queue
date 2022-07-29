package main

import (
	"log"
	"net/http"
	"sync"
)

func main() {
	log.Println("starting server")
	if err := runServer(":8080"); err != nil {
		log.Println(err)
	}
}

// runServer запускает HTTP сервер
func runServer(addr string) error {
	http.HandleFunc("/queue", getQueue)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getQueue(w, r)
		} else if r.Method == http.MethodPut {
			putQueue(w, r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}

	return nil
}

// getQueue отдает элемент из очереди
// где URL path - имя очереди
func getQueue(w http.ResponseWriter, r *http.Request) {
	queueName := r.URL.Path

	manager := ManagerInstance()
	if !manager.Exists(queueName) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	element, err := manager.Get(queueName).Get()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(element))
}

// putQueue добавляет данные в очередь
// где URL path - имя очереди, параметр v - добавляемые данные
func putQueue(w http.ResponseWriter, r *http.Request) {
	queueName := r.URL.Path
	element := r.URL.Query().Get("v")

	if queueName == "" || element == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ManagerInstance().Get(queueName).Add(element)
}

var m *Manager
var once sync.Once

// ManagerInstance возвращает менеджер очередей
func ManagerInstance() *Manager {
	once.Do(func() {
		m = NewManager()
	})

	return m
}
