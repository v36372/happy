package main

import (
	"net/http"

	"happy"
)

const songPerBatch = 20

func (a *App) IndexHandler(db *happy.PGDB) HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		songList, err := db.GetSongList(songPerBatch)
		if err != nil {
			a.logr.Log("error when getting song list: %s", err)
		}

		fs := a.getFlashes(w, req)

		pp := struct {
			*localPresenter
			SongList []*happy.Song
			Flashes  []interface{}
		}{
			localPresenter: &localPresenter{PageTitle: "", PageURL: "", globalPresenter: a.gp},
			SongList:       songList,
			Flashes:        fs,
		}
		err = a.rndr.HTML(w, http.StatusOK, "index", pp)
		if err != nil {
			a.logr.Log("error when rendering index page: %s", err)
		}
		return nil
	}
}

func (a *App) AboutHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		p := &localPresenter{PageTitle: "About grepbook", PageURL: "/about", globalPresenter: a.gp}
		err := a.rndr.HTML(w, http.StatusOK, "about", p)
		if err != nil {
			a.logr.Log("error when rendering about page: %s", err)
		}
		return nil
	}
}

func (a *App) NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	lp := &localPresenter{
		PageTitle:       "Page not found",
		PageURL:         "/404",
		globalPresenter: a.gp,
	}
	a.rndr.HTML(w, http.StatusNotFound, "404", lp)
	return
}
