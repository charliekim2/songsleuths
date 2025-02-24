package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
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
	} else if r.Method == http.MethodPatch {
		status, err = patch(w, r)
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

	// TODO: check num songs
	// Check if before deadline
	game := &db.Game{}
	res := conn.First(game, "id = ?", gid)
	if res.Error != nil {
		return http.StatusNotFound, res.Error
	}
	if time.Now().Unix() > int64(game.Deadline) {
		return http.StatusBadRequest, errors.New("deadline has passed")
	}

	var songs []db.Song
	idRegex := regexp.MustCompile(`[a-zA-Z0-9]{22}`)
	for _, song := range submission.Songs {
		if !idRegex.MatchString(song) {
			return http.StatusBadRequest, errors.New("invalid Spotify ID")
		}
		songs = append(songs, db.Song{
			Spotify: song,
			GameID:  gid,
		})
	}
	sub := db.Submission{
		PlayerID: uid,
		GameID:   gid,
		Nickname: submission.Nickname,
		Songs:    songs, // may need to do association append
		Drawing:  submission.Drawing,
	}
	err = conn.Create(&sub).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.WriteHeader(http.StatusCreated)
	return http.StatusCreated, nil
}

func remove(w http.ResponseWriter, r *http.Request) (int, error) {
	return http.StatusOK, nil
}

func patch(w http.ResponseWriter, r *http.Request) (int, error) {
	// TODO: consider doing upsert in post instead
	return http.StatusOK, nil
}
