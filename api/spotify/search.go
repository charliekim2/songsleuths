package handler

import (
	"errors"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusMethodNotAllowed
	err := errors.New("Invalid request method")

	if r.Method == http.MethodGet {
		status, err = get(w, r)
	}

	if err != nil {
		http.Error(w, err.Error(), status)
	}
}

func get(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusOK, nil
}
