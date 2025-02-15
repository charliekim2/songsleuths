package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/charliekim2/songsleuths/utils"
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

type SearchResponse struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Album   string   `json:"album"`
	Artists []string `json:"artists"`
	Image   string   `json:"image"`
}

func get(w http.ResponseWriter, r *http.Request) (int, error) {
	search := r.URL.Query().Get("q")
	if search == "" {
		return http.StatusBadRequest, errors.New("Search query is required")
	}

	spotifyResult, err := utils.Search(search, 10)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	var res []SearchResponse
	for _, track := range spotifyResult.Tracks.Items {
		artists := make([]string, len(track.Artists))
		for i, artist := range track.Artists {
			artists[i] = artist.Name
		}
		image := ""
		if len(track.Album.Images) > 0 {
			image = track.Album.Images[0].URL
		}
		res = append(res, SearchResponse{
			ID:      track.ID,
			Name:    track.Name,
			Album:   track.Album.Name,
			Artists: artists,
			Image:   image,
		})
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	return http.StatusOK, nil
}
