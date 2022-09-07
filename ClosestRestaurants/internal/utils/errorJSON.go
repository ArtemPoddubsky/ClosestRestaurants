package utils

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
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
