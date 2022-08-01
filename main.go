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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(ManagerInstance(), w, r)
	})

	log.Println("starting server on", addr)
	if err = http.ListenAndServe(addr, nil); err != nil {
		log.Println(err)
	}
}

// getAddr генерирует адрес для запуска сервера на основе аргументов программы
func getAddr() (string, error) {
	if len(os.Args) < 2 {
		return "", fmt.Errorf("port not specified, running: app.exe 8080")
	}

	return fmt.Sprintf(":%s", os.Args[1]), nil
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
