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

type Game struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Deadline uint   `json:"deadline"`
	NSongs   uint   `json:"n_songs"`

	// The requesting players submission
	Submission *Submission `json:"submission,omitempty"`
	// The other players who submitted
	PlayerList []struct {
		Nickname string `json:"nickname"`
		Drawing  string `json:"drawing"`
	} `json:"player_list,omitempty"`

	GuessList   *Tierlist `json:"guess_list,omitempty"`
	RankingList *Tierlist `json:"ranking_list,omitempty"`
	Playlist    string    `json:"playlist,omitempty"`
	Songs       []string  `json:"songs,omitempty"`
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
	// } else if r.Method == http.MethodPatch {
	// 	status, err = patch(w, r)
	// }

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
	res := conn.First(game, "id = ?", gid)
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
		for _, s := range game.Submissions {
			for _, songs := range s.Songs {
				g.Songs = append(g.Songs, songs.Spotify)
			}
		}
		g.GuessList = &Tierlist{}
		g.GuessList.ID = game.GuessList.ID
		g.GuessList.Type = game.GuessList.Type
		for _, tier := range game.GuessList.Tiers {
			g.GuessList.Tiers = append(g.GuessList.Tiers, Tier{
				ID:   tier.ID,
				Name: tier.Name,
				Rank: tier.Rank,
			})
		}
		g.RankingList = &Tierlist{}
		g.RankingList.ID = game.RankingList.ID
		g.RankingList.Type = game.RankingList.Type
		for _, tier := range game.RankingList.Tiers {
			g.RankingList.Tiers = append(g.RankingList.Tiers, Tier{
				ID:   tier.ID,
				Name: tier.Name,
				Rank: tier.Rank,
			})
		}
		g.Playlist = game.Playlist
	} else {
		sub := &db.Submission{}
		res = conn.Where(&db.Submission{PlayerID: uid, GameID: gid}).First(sub)
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
	// TODO: Check if player submitted rankings, and if so, return answers

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(g)
	return http.StatusOK, nil
}
