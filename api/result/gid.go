package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/charliekim2/songsleuths/db"
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

func get(w http.ResponseWriter, r *http.Request) (int, error) {
	uid, err := utils.Authenticate(r)
	if err != nil {
		return http.StatusUnauthorized, err
	}
	gid := strings.TrimPrefix(r.URL.Path, "/api/rank/")

	// Check if player submitted rankings for game
	conn, err := db.Connect()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	tierlist := &db.Tierlist{}
	err = conn.Where(&db.Tierlist{
		GameID: gid,
		Type:   "guess",
	}).First(tierlist).Error
	if err != nil {
		return http.StatusNotFound, err
	}
	var count int64
	if err := conn.Model(&db.Ranking{}).
		Where("player_id = ? AND tierlist_id != ?", uid, tierlist.ID).
		Count(&count).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	if count != 1 {
		return http.StatusBadRequest, errors.New("must submit guesses first")
	}

	w.WriteHeader(http.StatusOK)
	return 0, nil
}
