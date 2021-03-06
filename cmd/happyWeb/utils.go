package main

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

// saveFlash is a utility helper to save Flashes to the session cookie.
func (a *App) saveFlash(w http.ResponseWriter, req *http.Request, msg string) error {
	session, err := a.store.Get(req, SessionName)
	if err != nil {
		return err
	}
	session.AddFlash(msg)
	err = session.Save(req, w)
	if err != nil {
		return err
	}
	return nil
}

// getFlashes gets all the flases from request, and returns it.
func (a *App) getFlashes(w http.ResponseWriter, req *http.Request) []interface{} {
	session, _ := a.store.Get(req, SessionName)
	fs := session.Flashes()
	session.Save(req, w)
	return fs
}

// GetParamsObj returns a httprouter params object given the request.
func GetParamsObj(req *http.Request) httprouter.Params {
	ps := context.Get(req, Params).(httprouter.Params)
	return ps
}
