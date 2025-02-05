package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/charliekim2/songsleuths/db"
)

type game struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Deadline uint   `json:"deadline"`
	NSongs   uint   `json:"n_songs"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	status := 405
	err := errors.New("Invalid request method")

	if r.Method == http.MethodPost {
		status, err = post(w, r)
	}

	if err != nil {
		http.Error(w, err.Error(), status)
	}
}

func post(w http.ResponseWriter, r *http.Request) (int, error) {

	var g game
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		return 400, err
	}

	conn, err := db.Connect()
	if err != nil {
		return 500, err
	}

	dbGame := db.Game{
		Name:     g.Name,
		Deadline: g.Deadline,
		NSongs:   g.NSongs,
	}

	err = conn.Create(&dbGame).Error
	if err != nil {
		return 400, err
	}

	g.ID = dbGame.ID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	err = json.NewEncoder(w).Encode(g)
	if err != nil {
		return 500, err
	}
	return 201, nil
}
