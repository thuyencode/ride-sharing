package types

type Route struct {
	Distance float64     `json:"distance"`
	Duration float64     `json:"duration"`
	Geometry []*Geometry `json:"geometry"`
}

type Geometry struct {
	Coordinates []*Coordinate `json:"coordinates"`
}

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type OSRMAPIResponse struct {
	Code   string `json:"code"`
	Routes []struct {
		Legs []struct {
			Steps    []any   `json:"steps"`
			Weight   float64 `json:"weight"`
			Summary  string  `json:"summary"`
			Duration float64 `json:"duration"`
			Distance float64 `json:"distance"`
		} `json:"legs"`
		WeightName string `json:"weight_name"`
		Geometry   struct {
			Coordinates [][]float64 `json:"coordinates"`
			Type        string      `json:"type"`
		} `json:"geometry"`
		Weight   float64 `json:"weight"`
		Duration float64 `json:"duration"`
		Distance float64 `json:"distance"`
	} `json:"routes"`
	Waypoints []struct {
		Hint     string    `json:"hint"`
		Location []float64 `json:"location"`
		Name     string    `json:"name"`
		Distance float64   `json:"distance"`
	} `json:"waypoints"`
}
