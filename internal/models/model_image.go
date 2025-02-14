package models

type Image struct {
	Image    string   `json:"image"`
	Tag      string   `json:"tag"`
	Desc     string   `json:"desc"`
	Labels   []string `json:"labels"`
	Env      []string `json:"env"`
	Size     string   `json:"size"`
	PulledOn string   `json:"Pulled_on"`
}

type ImagesConfig struct {
	CustomImages []Image `json:"custom_images"`
}
