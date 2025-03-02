package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/charliekim2/songsleuths/db"
	"github.com/charliekim2/songsleuths/utils"
	"gorm.io/gorm/clause"
)

type Game struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Deadline uint   `json:"deadline"`
	NSongs   uint   `json:"n_songs"`

	// The requesting players submission
	Submission *Submission `json:"submission,omitempty"`
	// The other players who submitted
	// PlayerList []struct {
	// 	Nickname string `json:"nickname"`
	// 	Drawing  string `json:"drawing"`
	// } `json:"player_list,omitempty"`

	GuessList   *Tierlist `json:"guess_list,omitempty"`
	RankingList *Tierlist `json:"ranking_list,omitempty"`
	Playlist    string    `json:"playlist,omitempty"`
	Songs       []Song    `json:"songs,omitempty"`
}

type Song struct {
	ID       uint   `json:"id"`
	Spotify  string `json:"spotify"`
	AlbumArt string `json:"album_art"`
}

type Tierlist struct {
	ID    uint   `json:"id"`
	Type  string `json:"type"`
	Tiers []Tier `json:"tiers"`
}

type Tier struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Rank int    `json:"rank"`
}

type Submission struct {
	Songs    []string `json:"songs"`
	Nickname string   `json:"nickname"`
	Drawing  string   `json:"drawing"`
}

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

	gid := strings.TrimPrefix(r.URL.Path, "/api/games/")
	game := &db.Game{}
	conn, err := db.Connect()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	res := conn.Preload("Tierlists.Tiers").Preload("Submissions.Songs").Preload(clause.Associations).First(game, "id = ?", gid)
	if res.Error != nil {
		return http.StatusNotFound, res.Error
	}

	g := Game{
		ID:       game.ID,
		Name:     game.Name,
		Deadline: game.Deadline,
		NSongs:   game.NSongs,
	}
	if time.Now().Unix() > int64(game.Deadline) {
		// TODO: check addedSongs flag -> add songs, update flag
		g.Songs = []Song{}
		g.GuessList = &Tierlist{Tiers: []Tier{}}
		g.RankingList = &Tierlist{Tiers: []Tier{}}
		for _, s := range game.Submissions {
			for _, song := range s.Songs {
				g.Songs = append(g.Songs, Song{
					ID:       song.ID,
					Spotify:  song.Spotify,
					AlbumArt: song.AlbumArt,
				})
			}
		}
		if !game.AddedSongs {
			err = addSongs(g.Songs, game.Playlist)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			err = conn.Model(&db.Game{ID: game.ID}).Update("added_songs", true).Error
			if err != nil {
				return http.StatusInternalServerError, err
			}
		}
		for _, list := range game.Tierlists {
			if list.Type == "guess" {
				g.GuessList.ID = list.ID
				g.GuessList.Type = list.Type
				for _, tier := range list.Tiers {
					g.GuessList.Tiers = append(g.GuessList.Tiers, Tier{
						ID:   tier.ID,
						Name: tier.Name,
						Rank: tier.Rank,
					})
				}
			}
			if list.Type == "ranking" {
				g.RankingList.ID = list.ID
				g.RankingList.Type = list.Type
				for _, tier := range list.Tiers {
					g.RankingList.Tiers = append(g.RankingList.Tiers, Tier{
						ID:   tier.ID,
						Name: tier.Name,
						Rank: tier.Rank,
					})
				}
			}
		}
		g.Playlist = game.Playlist
	} else {
		sub := &db.Submission{}
		res = conn.Where(&db.Submission{PlayerID: uid, GameID: gid}).Preload("Songs").First(sub)
		if res.Error == nil {
			g.Submission = &Submission{
				Songs:    []string{},
				Nickname: sub.Nickname,
				Drawing:  sub.Drawing,
			}
			for _, s := range sub.Songs {
				g.Submission.Songs = append(g.Submission.Songs, s.Spotify)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(g)
	return http.StatusOK, nil
}

func addSongs(songs []Song, playlist string) error {
	songIds := []string{}
	for _, song := range songs {
		songIds = append(songIds, song.Spotify)
	}
	err := utils.AddToPlaylist(songIds, playlist)
	if err != nil {
		return err
	}
	return nil
}
