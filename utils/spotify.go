package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type tokenResponse struct {
	Token string `json:"access_token"`
}

func RefreshToken() (string, error) {
	auth_url := "https://accounts.spotify.com/api/token"
	auth := "Basic " + os.Getenv("BASE64_AUTH")
	content := "application/x-www-form-urlencoded"

	formData := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {os.Getenv("SPOTIFY_TOKEN")},
	}

	req, err := http.NewRequest("POST", auth_url, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", content)
	req.Header.Add("Authorization", auth)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errors.New("Spotify rejected token refresh request")
	}

	token := tokenResponse{}
	err = json.NewDecoder(res.Body).Decode(&token)
	if err != nil {
		return "", err
	}

	return token.Token, nil
}

func appToken() (string, error) {
	url := "https://accounts.spotify.com/api/token"
	content := "application/x-www-form-urlencoded"

	body := []byte("grant_type=client_credentials&client_id=" + os.Getenv("CLIENT_ID") + "&client_secret=" + os.Getenv("CLIENT_SECRET"))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", content)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errors.New("Spotify rejected token request")
	}

	token := tokenResponse{}
	err = json.NewDecoder(res.Body).Decode(&token)
	if err != nil {
		return "", err
	}

	return token.Token, nil
}

type SearchResult struct {
	Tracks struct {
		Items []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Album struct {
				Name   string `json:"name"`
				Images []struct {
					URL string `json:"url"`
				} `json:"images"`
			} `json:"album"`
		}
	}
}

func Search(query string, limit int) (*SearchResult, error) {
	token, err := appToken()
	if err != nil {
		return nil, err
	}

	url := "https://api.spotify.com/v1/search?q=" + url.QueryEscape(query) + "&type=track&limit=" + strconv.Itoa(limit)
	auth := "Bearer " + token

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", auth)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Spotify rejected search request")
	}

	result := SearchResult{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type uris struct {
	URIs []string `json:"uris"`
}

func AddToPlaylist(songs []string, playlist string) error {
	tok, err := RefreshToken()
	if err != nil {
		return err
	}

	addSongs := uris{
		URIs: []string{},
	}
	for _, songId := range songs {
		addSongs.URIs = append(addSongs.URIs, "spotify:track:"+songId)
	}
	body, err := json.Marshal(&addSongs)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "https://api.spotify.com/v1/playlists/"+playlist+"/tracks", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+tok)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return errors.New("Spotify rejected playlist creation request")
	}
	return nil
}

type TrackResult struct {
	Tracks []struct {
		Album struct {
			Images []struct {
				URL string `json:"url"`
			} `json:"images"`
		} `json:"album"`
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"tracks"`
}

type AlbumArt struct {
	ID   string
	URL  string
	Name string // show song name on hover
}

func GetAlbumArt(songs []string) ([]AlbumArt, error) {
	token, err := appToken()
	if err != nil {
		return nil, err
	}

	url := "https://api.spotify.com/v1/tracks?ids=" + strings.Join(songs, ",")
	auth := "Bearer " + token

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", auth)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Spotify rejected search request")
	}

	result := TrackResult{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	covers := []AlbumArt{}
	for _, t := range result.Tracks {
		if len(t.Album.Images) > 0 {
			covers = append(covers, AlbumArt{
				ID:   t.ID,
				URL:  t.Album.Images[0].URL,
				Name: t.Name,
			})
		} else {
			covers = append(covers, AlbumArt{
				ID:  t.ID,
				URL: "",
			})
		}
	}

	return covers, nil
}
