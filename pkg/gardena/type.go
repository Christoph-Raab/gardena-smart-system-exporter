package gardena

const (
	typeLocation = "LOCATION"
)

type LocationsFromApi struct {
	Data []struct {
		LocationFromApi
	} `json:"data"`
}

type LocationFromApi struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Name string `json:"name"`
	} `json:"attributes"`
}
