package main

import (
	"flag"
	"fmt"
	"happy"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
	"github.com/kardianos/osext"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/unrolled/render"
)

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
	port             string
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
			port:             viper.GetString("port"),
		}
	} else {
		cfg = config{
			youtubeAPIKey:    os.Getenv("youtubeAPIKey"),
			soundcloudAPIKey: os.Getenv("soundcloudAPIKey"),
			port:             os.Getenv("PORT"),
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
	devPath := pwd[:len(pwd)-3] + "cmd/happyWeb/"
	_, file, _, _ := runtime.Caller(1)
	configPath := path.Dir(file)
	viper.AddConfigPath(devPath)
	viper.AddConfigPath(configPath)
	return viper.ReadInConfig() // Find and read the config file
}

func init() {
	flag.StringVar(&dev, "dev", "true", "build for local dev")
}

func main() {
	flag.Parse()
	pwd, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatalf("cannot retrieve present working directory: %i", 0600, nil)
	}

	// Load configuration
	err = LoadConfiguration(pwd)
	if err != nil {
		panic(errors.Errorf("Fatal reading config file: %s \n", err))
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
	r.Get("/song", common.Then(a.Wrap(a.GetSongHandler(db))))
	r.Delete("/song/:id", common.Then(a.Wrap(a.DeleteSongHandler(db))))
	r.Get("/about", common.Then(a.Wrap(a.AboutHandler())))

	r.ServeFiles("/static/*filepath", http.Dir(staticFilePath))

	def := alice.New(responseWriterWrapper).Extend(common)
	r.NotFound = def.Then(responseWriterWrapper(http.HandlerFunc(a.NotFoundHandler)))

	err = http.ListenAndServe(":"+a.cfg.port, r)
	if err != nil {
		log.Println("error on serve server %s", err)
	}
}
