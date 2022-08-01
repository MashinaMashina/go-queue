package main

import (
	"errors"
	"net/http"
	"strings"
)

// handler - главный обработчик HTTP запросов
func handler(manager *Manager, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		res, err := getQueue(manager, strings.Trim(r.URL.Path, "/"))
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			w.Write([]byte(res))
		}
	} else if r.Method == http.MethodPut {
		err := putQueue(manager, strings.Trim(r.URL.Path, "/"), r.URL.Query().Get("v"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
