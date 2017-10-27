package main

import (
	"encoding/json"
	"fmt"
	"happy"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

const YT = "youtube"
const SC = "soundcloud"
const YT_LIST_VIDEO_API = "https://www.googleapis.com/youtube/v3/videos/"
const YT_PLAYLIST_API = "https://www.googleapis.com/youtube/v3/playlistItems?maxResults=50&part=snippet"
const SC_RESOLVE_API = "http://api.soundcloud.com/resolve/"

func (a *App) GetSongHandler(db *happy.PGDB) HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		songList, err := db.GetSongList()
		if err != nil {
			a.logr.Log("error when getting song list: %s", err)
		}

		err = json.NewEncoder(w).Encode(songList)
		if err != nil {
			a.logr.Log("error when return json %s", err)
			return newError(500, "error when return json %s", err)
		}

		return nil
	}
}

func (a *App) CreateSongHandler(db *happy.PGDB) HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		songLink := req.FormValue("link")

		if strings.TrimSpace(songLink) == "" {
			a.saveFlash(w, req, "Song address is invalid!")
			http.Redirect(w, req, "/", 302)
			return newError(400, "song link cannot be empty", nil)
		}

		var songs []*happy.Song
		var err error
		if strings.Index(songLink, YT) != -1 {
			u, err := url.Parse(songLink)
			if err != nil {
				a.logr.Log("error on CreateSongHandler, url parse fails: %s", err)
				return newError(500, "error on CreateSongHandler, url parse fails: %s", err)
			}

			params := u.Query()
			if u.Path == "/playlist" {
				songs, err = a.CreateSongFromYoutubePlaylist(params["list"][0])
			} else if u.Path == "/watch" {
				songs, err = a.CreateSongFromYoutube(params["v"][0])
			} else {
				a.logr.Log("unsupported link: %s", songLink)
				http.Redirect(w, req, "/", 302)
				return newError(400, "unsupported link: "+songLink, nil)
			}
			if err != nil {
				a.logr.Log("error on CreateSongHandler, CreateSongFromYoutube fails: %s", err)
				return newError(500, "error on CreateSongHandler, CreateSongFromYoutube fails: %s", err)
			}
		} else if strings.Index(songLink, SC) != -1 {
			songs, err = a.CreateSongFromSoundcloud(songLink)
			if err != nil {
				a.logr.Log("error on CreateSongHandler, CreateSongFromSoundcloud fails: %s", err)
				return newError(500, "error on CreateSongHandler, CreateSongFromSoundcloud fails: %s", err)
			}
		} else {
			a.logr.Log("unsupported music: %s", songLink)
			http.Redirect(w, req, "/", 302)
			return newError(400, "unsupported music: "+songLink, nil)
		}

		err = db.CreateSong(songs)
		if err != nil {
			a.logr.Log("error on creating song: %s", err)
			return newError(500, "error on creating song: %s", err)
		}
		http.Redirect(w, req, "/", 302)

		return nil
	}
}

func (a *App) CreateSongFromSoundcloud(link string) ([]*happy.Song, error) {
	resp, err := http.Get(SC_RESOLVE_API + fmt.Sprintf("?url=%s&client_id=%s", link, a.cfg.soundcloudAPIKey))

	if err != nil {
		return nil, err
	}

	scResolveResponse := happy.SCResolveResponse{}
	err = json.NewDecoder(resp.Body).Decode(&scResolveResponse)
	if err != nil {
		return nil, err
	}

	songs := []*happy.Song{}

	if scResolveResponse.Kind == "track" {
		newSong := &happy.Song{
			Name:      scResolveResponse.Title,
			Link:      scResolveResponse.URI,
			Provider:  "soundcloud",
			Thumbnail: scResolveResponse.ArtworkURL,
		}

		songs = append(songs, newSong)
	} else {
		for _, track := range scResolveResponse.Tracks {
			newSong := &happy.Song{
				Name:      track.Title,
				Link:      track.URI,
				Provider:  "soundcloud",
				Thumbnail: track.ArtworkURL,
			}
			songs = append(songs, newSong)
		}
	}
	return songs, nil
}

func (a *App) CreateSongFromYoutubePlaylist(playlistID string) ([]*happy.Song, error) {
	resp, err := http.Get(YT_PLAYLIST_API + fmt.Sprintf("&playlistId=%s&key=%s", playlistID, a.cfg.youtubeAPIKey))

	if err != nil {
		return nil, err
	}

	ytVideoListResponse := happy.YTVideoListResponse{}
	err = json.NewDecoder(resp.Body).Decode(&ytVideoListResponse)
	if err != nil {
		return nil, err
	}

	if len(ytVideoListResponse.Items) == 0 {
		return nil, errors.Errorf("youtube playlist with id %s not found", playlistID)
	}

	songs := []*happy.Song{}
	quality := []string{"maxres", "standard", "high", "medium", "default"}
	for _, song := range ytVideoListResponse.Items {
		thumbnails := song.Snippet.Thumbnails
		start := 0

		for start < len(quality)-1 {
			_, ok := thumbnails[quality[start]]
			if ok {
				break
			}
			start++
		}

		newSong := &happy.Song{
			Name:      song.Snippet.Title,
			Link:      song.Snippet.ResourceID.VideoID,
			Provider:  "youtube",
			Thumbnail: thumbnails[quality[start]].URL,
		}

		songs = append(songs, newSong)
	}

	return songs, nil
}

func (a *App) CreateSongFromYoutube(videoID string) ([]*happy.Song, error) {
	resp, err := http.Get(YT_LIST_VIDEO_API + fmt.Sprintf("?id=%s&part=snippet&key=%s", videoID, a.cfg.youtubeAPIKey))

	if err != nil {
		return nil, err
	}

	ytVideoListResponse := happy.YTVideoListResponse{}
	err = json.NewDecoder(resp.Body).Decode(&ytVideoListResponse)
	if err != nil {
		return nil, err
	}

	if len(ytVideoListResponse.Items) == 0 {
		return nil, errors.Errorf("youtube video with id %s not found", videoID)
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

	songs := []*happy.Song{}

	newSong := &happy.Song{
		Name:      ytVideoListResponse.Items[0].Snippet.Title,
		Link:      ytVideoListResponse.Items[0].ID,
		Provider:  "youtube",
		Thumbnail: thumbnails[quality[start]].URL,
	}

	songs = append(songs, newSong)
	return songs, nil
}
