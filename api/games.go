package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/charliekim2/songsleuths/db"
	"github.com/charliekim2/songsleuths/utils"
)

type game struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Deadline uint   `json:"deadline"`
	NSongs   uint   `json:"n_songs"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusMethodNotAllowed
	err := errors.New("Invalid request method")

	_, err = utils.Authenticate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

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
		return http.StatusBadRequest, err
	}

	conn, err := db.Connect()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	dbGame := db.Game{
		Name:     g.Name,
		Deadline: g.Deadline,
		NSongs:   g.NSongs,
	}

	err = conn.Create(&dbGame).Error
	if err != nil {
		return http.StatusBadRequest, err
	}

	g.ID = dbGame.ID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(g)
	return http.StatusCreated, nil
}
