package elasticsearch

type Restaurant struct {
	ID        string  `json: "ID"`
	Name      string  `json: "Name"`
	Address   string  `json: "Address"`
	Phone     string  `json: "Phone"`
	Longitude float64 `json: "Longitude"`
	Latitude  float64 `json: "Latitude"`
}

// next version //
/*
type Restaurant struct {
	ID        string  `json: "ID"`
	Name      string  `json: "Name"`
	Address   string  `json: "Address"`
	Phone     string  `json: "Phone"`
	type Location struct {
		Longitude float64 `json: "Longitude""
		Latitude float64 `json: "Latitude"`
	}
}
*/
