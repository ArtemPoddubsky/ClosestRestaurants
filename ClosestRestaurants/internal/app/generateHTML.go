package app

import (
	"html/template"
	"main/internal/utils"
	"net/http"
)

func GenerateHTML(w http.ResponseWriter, rests utils.HTMLPlaces) error {
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
	return t.Execute(w, rests)
}
