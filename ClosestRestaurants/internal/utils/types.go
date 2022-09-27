package utils

// HTMLPlaces is used to generate or get data in HTML.
type HTMLPlaces struct {
	Total  int
	Page   int
	Places []Restaurant
}

// JSONPlaces is used to form a JSON response consisting three closest restaurants.
type JSONPlaces struct {
	Name   string       `json:"name"`
	Places []Restaurant `json:"places"`
}

// Restaurant struct holds all data about one restaurant.
type Restaurant struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location Location `json:"location"`
}

// Location struct is exported to decode coordinates provided in request for API.
type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
