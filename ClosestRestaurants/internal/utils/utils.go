package utils

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func ErrorJSON(w http.ResponseWriter, msg string, code int) {
	v := struct {
		Error string `json:"error"`
	}{msg}
	w.WriteHeader(code)
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	if err := e.Encode(v); err != nil {
		http.Error(w, "Temporary Error in creating JSON response (500) ", http.StatusInternalServerError)
		log.Errorln("Couldn't Encode error response: ", err)
	}
}

func GracefullShutdown(server *http.Server) {
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-termChan
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalln("server.Shutdown:", err)
	}
	log.Infoln("Gracefully shutdown after signal:", sig)
}
