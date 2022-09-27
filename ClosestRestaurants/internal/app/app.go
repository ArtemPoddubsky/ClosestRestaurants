package app

import (
	"context"
	"encoding/json"
	"errors"
	"main/internal/config"
	"main/internal/storage"
	"main/internal/utils"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// App holds configuration data and storage instance.
type App struct {
	config *config.Config
	db     *storage.Postgres
}

// NewApp returns new instance of App.
func NewApp(cfg *config.Config) *App {
	return &App{
		config: cfg,
		db:     storage.NewDB(cfg),
	}
}

func (a *App) router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", a.List)
	router.HandleFunc("/api/recommend", a.APIRecommend)
	router.HandleFunc("/css/style.css", a.serveStyle)

	return router
}

// Run fills database if needed, configures and runs server with GracefullShutdown.
func (a *App) Run() {
	if err := a.db.FillDatabase(); err != nil {
		log.Fatalln("FillDatabase: ", err)
	}

	server := &http.Server{
		Addr:         ":5000",
		Handler:      a.router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Infoln("Running")

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("ListenAndServe:", err)
		}
	}()

	utils.GracefullShutdown(server)
}

func (a *App) serveStyle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./materials/css/style.css")
}

// List parses requested page number and responds with generated HTML page with list of restaurants.
func (a *App) List(w http.ResponseWriter, r *http.Request) {
	values, err := url.ParseQuery(r.RequestURI)
	if err != nil {
		http.Error(w, "Temporary Error (500)", http.StatusInternalServerError)
		log.Errorln("List: url.ParseQuery:", err)

		return
	}

	page, err := strconv.Atoi(values.Get("/?page"))
	if err != nil || page < 0 {
		http.Error(w, "This page can't possibly exist", http.StatusBadRequest)

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	rests, err := a.db.GetPage(ctx, page)

	cancel()

	if err != nil {
		if err.Error() == "Not Found" {
			http.Error(w, "This page doesn't exist", http.StatusNotFound)
		} else {
			http.Error(w, "Temporary Error (500) ", http.StatusInternalServerError)
			log.Errorln("List: db.GetPage: ", err, r.RequestURI)
		}

		return
	}

	if err = generateHTML(w, rests); err != nil {
		http.Error(w, "Temporary Error (500) ", http.StatusInternalServerError)
		log.Errorln("List: GenerateHTML: ", err)

		return
	}
}

// APIRecommend responds with JSON representation of 3 closest restaurants
// based on location sent in request.
func (a *App) APIRecommend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var coor utils.Location
	err := json.NewDecoder(r.Body).Decode(&coor)

	if err != nil || (coor.Lat == 0 || coor.Lon == 0) {
		utils.ErrorJSON(w, "Bad JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	rests, err := a.db.GetClosest(ctx, coor.Lat, coor.Lon)

	cancel()

	if err != nil {
		log.Errorln("API: db.ThreeClosest: Temporary Error (500) ", err)
		utils.ErrorJSON(w, "Temporary Error (500)", http.StatusInternalServerError)

		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err = encoder.Encode(rests); err != nil {
		log.Errorln("API: ThreeClosest: Temporary Error (500): Encode ", err)
		utils.ErrorJSON(w, "Temporary Error (500)", http.StatusInternalServerError)

		return
	}
}
