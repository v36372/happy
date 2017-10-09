package main

import (
	"flag"
	"fmt"
	"happy"
	"net/http"
	"path"

	"github.com/gorilla/context"
	"github.com/justinas/alice"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/unrolled/render"
)

var configPath string

type App struct {
	rndr   *render.Render
	router *Router
	gp     globalPresenter
	bm     *bluemonday.Policy
	logr   appLogger
}

// localPresenter contains the fields nessessary for specific pages.
type localPresenter struct {
	PageTitle string
	PageURL   string
	globalPresenter
}

type globalPresenter struct {
	SiteName    string
	Description string
	SiteURL     string
}

func SetupApp(r *Router, logger appLogger, templateFolderPath string) *App {
	rndr := render.New(render.Options{
		Directory:  templateFolderPath,
		Extensions: []string{".html"},
		Layout:     "base",
	})

	gp := globalPresenter{
		SiteName:    "Happy",
		Description: "Music player for happy people",
		SiteURL:     "happy.laptrinhviendeptrai.xyz",
	}

	bm := bluemonday.UGCPolicy()

	return &App{
		rndr:   rndr,
		router: r,
		gp:     gp,
		bm:     bm,
		logr:   logger,
	}
}

func LoadConfiguration(pwd string) error {
	viper.SetConfigName("happy-config")
	viper.AddConfigPath(pwd)

	return viper.ReadInConfig()
}

func init() {
	flag.StringVar(&configPath, "config", "", "path to config file for this app")
}

func main() {
	flag.Parse()

	// Load configuration
	err := LoadConfiguration(configPath)
	if err != nil {
		panic(errors.Errorf("Fatal reading config file: %s \n", err))
	}

	// Set up application path

	staticFilePath := path.Join(viper.GetString("path"), "static")
	templateFolderPath := path.Join(viper.GetString("path"), "templates")

	// Set up Database
	db, err := happy.OpenDB(fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		viper.GetString("databaseURL"),
		viper.GetInt("databasePort"),
		viper.GetString("databaseUser"),
		viper.GetString("databasePass"),
		viper.GetString("databaseName")))
	if err != nil {
		panic(errors.Errorf("Cannot connect to database: %s", err))
	}

	r := NewRouter()
	logr := newLogger()
	a := SetupApp(r, logr, templateFolderPath)

	common := alice.New(context.ClearHandler, a.loggingHandler, a.recoverHandler)

	r.Get("/", common.Then(a.Wrap(a.IndexHandler(db))))
	r.ServeFiles("/static/*filepath", http.Dir(staticFilePath))

	def := alice.New(responseWriterWrapper).Extend(common)
	r.NotFound = def.Then(responseWriterWrapper(http.HandlerFunc(a.NotFoundHandler)))

	http.ListenAndServe(":3000", r)
}
