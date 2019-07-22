package geo

// City becouse geoip2.City is a big inside our data
type City struct {
	City        string
	Postal      string
	Subdivision struct { // AKA state, province etc
		Name string
		ISO  string
	}
	Continent string
	Country   string
	Location  struct {
		AccuracyRadius uint16
		Latitude       float64
		Longitude      float64
		TimeZone       string
	}
	Traits struct { // Extra shit
		IsProxy bool
		IsEU    bool
	}
}
