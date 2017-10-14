package main

import (
	"encoding/json"
	"fmt"
	"happy"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const YT = "youtube"
const SC = "soundcloud"
const YT_LIST_VIDEO_API = "https://www.googleapis.com/youtube/v3/videos/"
const SC_RESOLVE_API = "http://api.soundcloud.com/resolve/"

func (a *App) CreateSongHandler(db *happy.PGDB) HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		songLink := req.FormValue("link")

		if strings.TrimSpace(songLink) == "" {
			a.saveFlash(w, req, "Song address is invalid!")
			http.Redirect(w, req, "/", 302)
			return newError(400, "song link cannot be empty", nil)
		}

		youtubeID := songLink[strings.Index(songLink, "=")+1:]

		var song *happy.Song
		var err error
		if strings.Index(songLink, YT) != -1 {
			song, err = a.CreateSongFromYoutube(youtubeID)
			if err != nil {
				a.logr.Log("error on CreateSongHandler, CreateSongFromYoutube fails: %s", err)
				return newError(500, "error on CreateSongHandler, CreateSongFromYoutube fails: %s", err)
			}
		} else if strings.Index(songLink, SC) != -1 {
			song, err = a.CreateSongFromSoundcloud(songLink)
			if err != nil {
				a.logr.Log("error on CreateSongHandler, CreateSongFromSoundcloud fails: %s", err)
				return newError(500, "error on CreateSongHandler, CreateSongFromSoundcloud fails: %s", err)
			}
		} else {
			a.logr.Log("unsupported music: %s", songLink)
			http.Redirect(w, req, "/", 302)
			return newError(400, "unsupported music: "+songLink, nil)
		}

		err = db.CreateSong(song)
		if err != nil {
			a.logr.Log("error on creating song: %s", err)
			return newError(500, "error on creating song: %s", err)
		}
		http.Redirect(w, req, "/", 302)

		return nil
	}
}

func (a *App) CreateSongFromSoundcloud(link string) (*happy.Song, error) {
	resp, err := http.Get(SC_RESOLVE_API + fmt.Sprintf("?url=%s&client_id=%s", link, a.cfg.soundcloudAPIKey))

	if err != nil {
		return nil, err
	}

	scResolveResponse := happy.SCResolveResponse{}
	err = json.NewDecoder(resp.Body).Decode(&scResolveResponse)
	if err != nil {
		return nil, err
	}

	return &happy.Song{
		Name:      scResolveResponse.Title,
		Link:      scResolveResponse.URI,
		Provider:  "soundcloud",
		Thumbnail: scResolveResponse.ArtworkURL,
	}, nil
}

func (a *App) CreateSongFromYoutube(link string) (*happy.Song, error) {
	resp, err := http.Get(YT_LIST_VIDEO_API + fmt.Sprintf("?id=%s&part=snippet&key=%s", link, a.cfg.youtubeAPIKey))

	if err != nil {
		return nil, err
	}

	ytVideoListResponse := happy.YTVideoListResponse{}
	err = json.NewDecoder(resp.Body).Decode(&ytVideoListResponse)
	if err != nil {
		return nil, err
	}

	if len(ytVideoListResponse.Items) == 0 {
		return nil, errors.Errorf("youtube video with id %s not found", link)
	}

	quality := []string{"maxres", "standard", "high", "medium", "default"}
	thumbnails := ytVideoListResponse.Items[0].Snippet.Thumbnails
	start := 0

	for start < len(quality)-1 {
		_, ok := thumbnails[quality[start]]
		if ok {
			break
		}
		start++
	}

	return &happy.Song{
		Name:      ytVideoListResponse.Items[0].Snippet.Title,
		Link:      ytVideoListResponse.Items[0].ID,
		Provider:  "youtube",
		Thumbnail: thumbnails[quality[start]].URL,
	}, nil
}
