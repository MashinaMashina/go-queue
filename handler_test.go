package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(method string, manager *Manager, queue, element string) (int, string, error) {
	r := httptest.NewRequest(method, "/"+queue+"?v="+element, nil)
	w := httptest.NewRecorder()

	handler(manager, w, r)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, "", fmt.Errorf("error with reading put response. Error: %w", err)
	}

	return res.StatusCode, string(data), nil
}

func testRequestPut(manager *Manager, queue, element string) (int, string, error) {
	code, resp, err := testRequest(http.MethodPut, manager, queue, element)

	if resp != "" {
		return 0, "", fmt.Errorf("put response not empty")
	}

	return code, resp, err
}

func testRequestGet(manager *Manager, queue string) (int, string, error) {
	return testRequest(http.MethodGet, manager, queue, "")
}

func TestHandlerPutEmpty(t *testing.T) {
	queueName := "test-queue"
	manager := NewManager()

	code, resp, err := testRequestPut(manager, queueName, "")

	if err != nil {
		t.Error(err)
	}

	if resp != "" {
		t.Errorf("put response not empty")
	}

	if code != http.StatusBadRequest {
		t.Errorf("put status code not expected. Extected %d!= real %d", http.StatusBadRequest, code)
	}
}

func TestHandlerPut(t *testing.T) {
	queueName := "test-queue"
	element := "test-element"
	manager := NewManager()

	code, resp, err := testRequestPut(manager, queueName, element)

	if err != nil {
		t.Error(err)
	}

	if resp != "" {
		t.Errorf("put response not empty")
	}

	if code != http.StatusOK {
		t.Errorf("put status code not expected. Extected %d!= real %d", http.StatusOK, code)
	}

	res, err := manager.Get(queueName).Get()
	if err != nil {
		t.Errorf("getting wrote data from queue. Got error: %s", err.Error())
	}

	if res != element {
		t.Errorf("element in queue not expected: real '%s' != expected '%s'", res, element)
	}
}

func TestGetFIFO(t *testing.T) {
	queueName := "test-queue"
	manager := NewManager()
	elements := []string{
		"test-element-1",
		"test-element-2",
		"test-element-3",
	}

	queue := manager.Get(queueName)
	for _, element := range elements {
		queue.Add(element)
	}

	for i := len(elements) - 1; i >= 0; i-- {
		code, resp, err := testRequestGet(manager, queueName)
		if err != nil {
			t.Errorf("not expected error '%s'", err.Error())
		}
		if code != http.StatusOK {
			t.Errorf("extected http status %d, but real %d", http.StatusOK, code)
		}
		if elements[i] != resp {
			t.Errorf("expected from queue '%s', but real '%s'", elements[i], resp)
		}
	}

	code, _, _ := testRequestGet(manager, queueName)
	if code != http.StatusNotFound {
		t.Errorf("expected code %d with empty queue, but real %d", http.StatusNotFound, code)
	}
}
