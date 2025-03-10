package handler

import (
	"encoding/json"
	"errors"
	"fmt"
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

type Submission struct {
	Nickname string   `json:"nickname"`
	Songs    []string `json:"songs"`
	Drawing  string   `json:"drawing"`
}

func post(w http.ResponseWriter, r *http.Request) (int, error) {
	uid, err := utils.Authenticate(r)
	if err != nil {
		return http.StatusUnauthorized, err
	}
	gid := strings.TrimPrefix(r.URL.Path, "/api/submit/")
	submission := &Submission{}
	err = json.NewDecoder(r.Body).Decode(submission)
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
	if time.Now().Unix() > int64(game.Deadline) {
		return http.StatusBadRequest, errors.New("deadline has passed")
	}
	if len(submission.Songs) != int(game.NSongs) {
		return http.StatusBadRequest, errors.New(fmt.Sprintf("number of songs should be %d", game.NSongs))
	}

	err = conn.Unscoped().Where("player_id = ? and game_id = ?", uid, gid).Delete(&db.Submission{}).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}

	var songs []db.Song
	covers, err := utils.GetAlbumArt(submission.Songs)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	for _, cover := range covers {
		songs = append(songs, db.Song{
			Spotify:  cover.ID,
			AlbumArt: cover.URL,
			Name:     cover.Name,
			GameID:   gid,
		})
	}
	sub := db.Submission{
		PlayerID: uid,
		GameID:   gid,
		Nickname: submission.Nickname,
		Songs:    songs,
		Drawing:  submission.Drawing,
	}
	err = conn.Create(&sub).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.WriteHeader(http.StatusCreated)
	return 0, nil
}

func remove(w http.ResponseWriter, r *http.Request) (int, error) {
	uid, err := utils.Authenticate(r)
	if err != nil {
		return http.StatusUnauthorized, err
	}
	gid := strings.TrimPrefix(r.URL.Path, "/api/submit/")

	conn, err := db.Connect()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	game := &db.Game{}
	err = conn.First(game, "id = ?", gid).Error
	if err != nil {
		return http.StatusNotFound, err
	}
	if time.Now().Unix() > int64(game.Deadline) {
		return http.StatusBadRequest, errors.New("deadline has passed")
	}

	err = conn.Unscoped().Where("player_id = ? and game_id = ?", uid, gid).Delete(&db.Submission{}).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.WriteHeader(http.StatusNoContent)
	return 0, nil
}
