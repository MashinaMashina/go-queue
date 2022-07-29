package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	addr, err := getAddr()

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("starting server on", addr)
	if err = runServer(addr); err != nil {
		log.Println(err)
	}
}

func getAddr() (string, error) {
	if len(os.Args) < 2 {
		return "", fmt.Errorf("port not specified, running: app.exe 8080")
	}

	return fmt.Sprintf(":%s", os.Args[1]), nil
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
