package models

type Restaurant struct {
	ID       string `json: "ID"`
	Name     string `json: "Name"`
	Address  string `json: "Address"`
	Phone    string `json: "Phone"`
	Location struct {
		Longitude float64 `json: "Longitude"`
		Latitude  float64 `json: "Latitude"`
	} `json: "Location"`
}
