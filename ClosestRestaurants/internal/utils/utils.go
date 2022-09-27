package utils

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// ErrorJSON replies to the request with the specified error message in JSON format and HTTP code.
func ErrorJSON(writer http.ResponseWriter, msg string, code int) {
	writer.WriteHeader(code)

	response := struct {
		Error string `json:"error"`
	}{msg}
	e := json.NewEncoder(writer)
	e.SetIndent("", "  ")

	if err := e.Encode(response); err != nil {
		http.Error(writer, "Temporary Error in creating JSON response (500) ", http.StatusInternalServerError)
		log.Errorln("Couldn't Encode JSON error response: ", err)
	}
}

// GracefullShutdown is waiting for SIGINT or SIGTERM signals from http.Server to perform shutdown and log this action.
func GracefullShutdown(server *http.Server) {
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-termChan

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalln("server.Shutdown:", err)
	}

	log.Infoln("Gracefully shutdown after signal:", sig)
}
