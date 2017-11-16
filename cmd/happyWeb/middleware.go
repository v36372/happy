package main

import (
	"happy"
	"log"
	"net/http"
	"os"
	"time"
)

const SessionName = "session-happy-2842732"

// AppLogger is an interface for logging
type appLogger interface {
	Log(str string, v ...interface{})
}

type happyLogger struct {
	*log.Logger
}

// Log products a log entry with the current time prepended
func (hl *happyLogger) Log(str string, v ...interface{}) {
	v = append(v, 0)
	copy(v[1:], v[0:])
	v[0] = happy.TimeNow().Format(time.RFC3339)

	hl.Printf("[%s] "+str, v...)
}

func newLogger() *happyLogger {
	return &happyLogger{
		log.New(os.Stdout, "[happy] ", 0),
	}
}

// loggingHandler middleware logs all request
func (a *App) corsMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "OPTIONS" {
			http.StatusText(204)
		}
	}
	return http.HandlerFunc(fn)
}

// loggingHandler middleware logs all request
func (a *App) loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		t1 := happy.TimeNow()
		a.logr.Log("Started %s %s", req.Method, req.URL.Path)

		next.ServeHTTP(w, req)

		rw, ok := w.(ResponseWriter)
		if ok {
			a.logr.Log("Completed %v %s in %v", rw.Status(), http.StatusText(rw.Status()), time.Since(t1))
		} else {
			a.logr.Log("Unable to log due to invalid ResponseWriter conversion")
		}
	}
	return http.HandlerFunc(fn)
}

// responseWriterWrapper is a middleware to ensure that a http.Handler
// uses our own custom ResponseWriter, in order to capture response data
func responseWriterWrapper(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(NewResponseWriter(w), req)
	}
	return http.HandlerFunc(fn)
}

// recoverHandler is a middleware that captures and recovers from panics.
func (a *App) recoverHandler(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				a.logr.Log("Panic: %+v", err)
				http.Error(rw, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(rw, req)
	}
	return http.HandlerFunc(fn)
}
