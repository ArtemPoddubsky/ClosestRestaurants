package app

import (
	"fmt"
	"html/template"
	"main/internal/utils"
	"net/http"
)

func generateHTML(w http.ResponseWriter, rests utils.HTMLPlaces) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	funcTemp := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"dec": func(i int) int {
			return i - 1
		},
		"div": func(a, b int) int {
			return a / b
		},
	}
	t := template.Must(template.New("page.html").Funcs(funcTemp).ParseFiles("./materials/page.html"))

	if err := t.Execute(w, rests); err != nil {
		return fmt.Errorf("template.Execute: %w", err)
	}

	return nil
}
