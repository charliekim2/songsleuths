package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/charliekim2/songsleuths/db"
	"github.com/charliekim2/songsleuths/utils"
)

type game struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Deadline uint   `json:"deadline"`
	NSongs   uint   `json:"n_songs"`
}

type playlistRequest struct {
	Name        string `json:"name"`
	Public      bool   `json:"public"`
	Description string `json:"description"`
}

type playlistResponse struct {
	ID string `json:"id"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusMethodNotAllowed
	err := errors.New("Invalid request method")

	if r.Method == http.MethodPost {
		status, err = post(w, r)
	}

	if err != nil {
		http.Error(w, err.Error(), status)
	}
}

func post(w http.ResponseWriter, r *http.Request) (int, error) {
	_, err := utils.Authenticate(r)
	if err != nil {
		return http.StatusUnauthorized, err
	}

	var g game
	err = json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		return http.StatusBadRequest, err
	}

	conn, err := db.Connect()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// TODO: Move this logic to spotify utility
	tok, err := utils.RefreshToken()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	playlist := playlistRequest{
		Name:        g.Name,
		Description: "Song Sleuths playlist for " + g.Name,
		Public:      false,
	}
	body, err := json.Marshal(&playlist)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	req, err := http.NewRequest("POST", "https://api.spotify.com/v1/users/charliekim451/playlists", bytes.NewBuffer(body))
	if err != nil {
		return http.StatusInternalServerError, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+tok)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return res.StatusCode, errors.New("Spotify rejected playlist creation request")
	}
	playlistId := playlistResponse{}
	err = json.NewDecoder(res.Body).Decode(&playlistId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	dbGame := db.Game{
		Name:     g.Name,
		Deadline: g.Deadline,
		NSongs:   g.NSongs,
		Playlist: playlistId.ID,
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
