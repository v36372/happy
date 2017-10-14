package main

import (
	"flag"
	"fmt"
	"happy"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/unrolled/render"
)

var configPath string
var dev string

type App struct {
	rndr   *render.Render
	store  *sessions.CookieStore
	router *Router
	gp     globalPresenter
	cfg    config
	bm     *bluemonday.Policy
	logr   appLogger
}

type config struct {
	youtubeAPIKey    string
	soundcloudAPIKey string
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

func SetupApp(r *Router, logger appLogger, cookieSecret []byte, templateFolderPath string) *App {
	rndr := render.New(render.Options{
		Directory:  templateFolderPath,
		Extensions: []string{".html"},
	})

	cfg := config{}

	if dev == "true" {
		cfg = config{
			youtubeAPIKey:    viper.GetString("youtubeAPIKey"),
			soundcloudAPIKey: viper.GetString("soundcloudAPIKey"),
		}
	} else {
		cfg = config{
			youtubeAPIKey:    os.Getenv("youtubeAPIKey"),
			soundcloudAPIKey: os.Getenv("soundcloudAPIKey"),
		}
	}

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
		cfg:    cfg,
		bm:     bm,
		logr:   logger,
		store:  sessions.NewCookieStore(cookieSecret),
	}
}

func LoadConfiguration(pwd string) error {
	viper.SetConfigName("happy-config")
	viper.AddConfigPath(pwd)

	return viper.ReadInConfig() // Find and read the config file
}

func init() {
	flag.StringVar(&configPath, "config", "", "path to config file for this app")
	flag.StringVar(&dev, "dev", "", "build for local dev")
}

func main() {
	flag.Parse()

	// Load configuration
	var err error
	if dev == "true" {
		err = LoadConfiguration(configPath)
		if err != nil {
			panic(errors.Errorf("Fatal reading config file: %s \n", err))
		}
	}

	var dbURL, dbUser, dbPass, dbName, cookieSecret, appPath string
	var dbPort int

	if dev == "true" {
		dbURL = viper.GetString("databaseURL")
		dbPort = viper.GetInt("databasePort")
		dbUser = viper.GetString("databaseUser")
		dbPass = viper.GetString("databasePass")
		dbName = viper.GetString("databaseName")
		cookieSecret = viper.GetString("cookieSecret")
		appPath = viper.GetString("path")
	} else {
		dbURL = os.Getenv("DATABASE_URL")
		cookieSecret = os.Getenv("cookieSecret")
		appPath = os.Getenv("path")
	}
	// Set up application path

	staticFilePath := path.Join(appPath, "static")
	templateFolderPath := path.Join(appPath, "templates")

	// Set up Database
	var db *happy.PGDB
	if dev == "true" {
		db, err = happy.OpenDB(fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			dbURL, dbPort, dbUser, dbPass, dbName))
		if err != nil {
			panic(errors.Errorf("Cannot connect to database: %s", err))
		}
	} else {
		db, err = happy.OpenDB(dbURL)
		if err != nil {
			panic(errors.Errorf("Cannot connect to database: %s", err))
		}
	}
	r := NewRouter()
	logr := newLogger()

	a := SetupApp(r, logr, []byte(cookieSecret), templateFolderPath)

	common := alice.New(context.ClearHandler, a.loggingHandler, a.recoverHandler)

	r.Get("/", common.Then(a.Wrap(a.IndexHandler(db))))

	r.Post("/song", common.Then(a.Wrap(a.CreateSongHandler(db))))
	r.Get("/about", common.Then(a.Wrap(a.AboutHandler())))

	r.ServeFiles("/static/*filepath", http.Dir(staticFilePath))

	def := alice.New(responseWriterWrapper).Extend(common)
	r.NotFound = def.Then(responseWriterWrapper(http.HandlerFunc(a.NotFoundHandler)))

	http.ListenAndServe(":3000", r)
}
