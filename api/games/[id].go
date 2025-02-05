package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/charliekim2/songsleuths/db"
)

type game struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Deadline      uint   `json:"deadline"`
	NSongs        uint   `json:"n_songs"`
	GuessListID   uint   `json:"guess_list_id"`
	RankingListID uint   `json:"ranking_list_id"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	status := 405
	err := errors.New("Invalid request method")

	// if r.Method == http.MethodGet {
	// 	status, err = get(w, r)
	// } else if r.Method == http.MethodPatch {
	// 	status, err = patch(w, r)
	// }

	if err != nil {
		http.Error(w, err.Error(), status)
	}
}
