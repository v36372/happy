package main

import (
	"net/http"

	"happy"
)

func (a *App) IndexHandler(db *happy.PGDB) HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		pp := struct {
			*localPresenter
		}{
			localPresenter: &localPresenter{PageTitle: "Happy Music Player", PageURL: "/", globalPresenter: a.gp},
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
