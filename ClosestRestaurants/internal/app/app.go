package app

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"main/internal/config"
	"main/internal/storage"
	"main/internal/utils"
	"net/http"
	"net/url"
	"strconv"
)

type App struct {
	config config.Config
	db     *storage.Postgres
}

func NewApp(cfg config.Config) *App {
	return &App{
		config: cfg,
		db:     storage.NewDB(cfg),
	}
}

func (a *App) Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", a.List)
	router.HandleFunc("/css/style.css", a.ServeStyle)
	router.HandleFunc("/api/recommend", a.ApiRecommend)
	return router
}

func (a *App) Run() {
	if err := a.db.FillDatabase(); err != nil {
		log.Fatalln("Filling database error: ", err)
	}
	log.Infoln("Running")
	server := &http.Server{
		Addr:    "127.0.0.1:5000",
		Handler: a.Router(),
	}
	go func() {
		if servErr := server.ListenAndServe(); servErr != http.ErrServerClosed {
			log.Fatalln("ListenAndServe:", servErr)
		}
	}()
	utils.GracefullShutdown(server)
}

func (a *App) List(w http.ResponseWriter, r *http.Request) {
	u, err := url.ParseQuery(r.RequestURI)
	if err != nil {
		http.Error(w, "Temporary Error (500)", http.StatusInternalServerError)
		log.Errorln("url.ParseQuery:", err)
		return
	}
	page, err := strconv.Atoi(u.Get("/?page"))
	if err != nil || page < 0 {
		http.Error(w, "This page can't possibly exist", http.StatusBadRequest)
		return
	}

	rests, err := a.db.GetPage(page)
	if err != nil {
		if err.Error() == "Not Found" {
			http.Error(w, "This page doesn't exist", http.StatusNotFound)
		} else {
			http.Error(w, "Temporary Error (500) ", http.StatusInternalServerError)
			log.Errorln("Db.GetPage: ", err, r.RequestURI)
		}
		return
	}

	if err = GenerateHTML(w, rests); err != nil {
		http.Error(w, "Temporary Error (500) ", http.StatusInternalServerError)
		log.Errorln("GenerateHTML: ", err)
		return
	}
}

func (a *App) ServeStyle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./materials/css/style.css")
}

func (a *App) ApiRecommend(w http.ResponseWriter, r *http.Request) {
	var coor utils.Location
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&coor); err != nil || (coor.Lat == 0 || coor.Lon == 0) {
		utils.ErrorJSON(w, "Bad JSON", http.StatusBadRequest)
		return
	}
	rests, err := a.db.ThreeClosest(coor.Lat, coor.Lon)
	if err != nil {
		log.Errorln("Temporary Error (500) ", err)
		utils.ErrorJSON(w, "Temporary Error (500)", http.StatusInternalServerError)
		return
	}
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	if err = e.Encode(rests); err != nil {
		log.Errorln("Temporary Error (500) ", err)
		utils.ErrorJSON(w, "Temporary Error (500)", http.StatusInternalServerError)
		return
	}
}
