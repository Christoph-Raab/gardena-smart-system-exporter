package gardena

const (
	typeLocation = "LOCATION"
)

type Locations struct {
	Data []struct {
		Location
	} `json:"data"`
}

type Location struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Attribute
	} `json:"attributes"`
}

type Device struct {
	Id           string               `json:"id"`
	Type         string               `json:"type"`
	Relationship string               `json:"relationship"`
	Attributes   map[string]Attribute `json:"attributes"`
}

type Attribute struct {
	Name  string
	Value any
}

type State struct {
	Data struct {
		Id            string `json:"id"`
		Type          string `json:"type"`
		Relationships struct {
			Devices struct {
				Data []struct {
					Device
				} `json:"data"`
			} `json:"devices"`
		} `json:"relationships"`
		Attributes struct {
			Attribute
		} `json:"attributes"`
	} `json:"data"`
	Included []Device `json:"included"`
}
