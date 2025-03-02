package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/charliekim2/songsleuths/db"
	"github.com/charliekim2/songsleuths/utils"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	status, err := http.StatusMethodNotAllowed, errors.New("Invalid request method")

	if r.Method == http.MethodPost {
		status, err = post(w, r)
	} else if r.Method == http.MethodDelete {
		status, err = remove(w, r)
	}

	if err != nil {
		http.Error(w, err.Error(), status)
	}
}

type Ranking struct {
	TierlistID uint   `json:"tierlist_id"`
	Ranking    string `json:"ranking"`
}

func post(w http.ResponseWriter, r *http.Request) (int, error) {
	uid, err := utils.Authenticate(r)
	if err != nil {
		return http.StatusUnauthorized, err
	}
	gid := strings.TrimPrefix(r.URL.Path, "/api/rank/")
	ranking := &Ranking{}
	err = json.NewDecoder(r.Body).Decode(ranking)
	if err != nil {
		return http.StatusBadRequest, err
	}

	conn, err := db.Connect()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	game := &db.Game{}
	err = conn.First(game, "id = ?", gid).Error
	if err != nil {
		return http.StatusNotFound, err
	}
	if time.Now().Unix() < int64(game.Deadline) {
		return http.StatusBadRequest, errors.New("cannot rank before deadline")
	}

	// Cannot resubmit ranking
	// err = conn.Where("player_id = ? and tierlist_id = ?", uid, ranking.TierlistID).Delete(&db.Ranking{}).Error
	// if err != nil {
	// 	return http.StatusInternalServerError, err
	// }

	dbRanking := db.Ranking{
		PlayerID:   uid,
		TierlistID: ranking.TierlistID,
		GameID:     gid,
		Ranking:    ranking.Ranking,
	}
	err = conn.Create(&dbRanking).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.WriteHeader(http.StatusCreated)
	return 0, nil
}
