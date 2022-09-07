package utils

type HTMLPlaces struct {
	Total  int
	Page   int
	Places []Restaurant
}

type JSONPlaces struct {
	Name   string       `json:"name"`
	Places []Restaurant `json:"places"`
}

type Restaurant struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location Location `json:"location"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
