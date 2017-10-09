package main

import (
	"happy"
	"net/http"
	"strings"
)

func (a *App) CreateSongHandler(db *happy.PGDB) HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		songLink := req.FormValue("link")

		if strings.TrimSpace(songLink) == "" {
			a.saveFlash(w, req, "Song address is invalid!")
			http.Redirect(w, req, "/", 302)
			return newError(400, "song link cannot be empty", nil)
		}

		err := db.CreateSong(songLink)
		if err != nil {
			a.logr.Log("error on creating song: %s", err)
			return err
		}
		http.Redirect(w, req, "/", 302)

		return nil
	}
}
