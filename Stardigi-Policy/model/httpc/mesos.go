package httpc

type Mesoses struct {
	Frameworks []Framwork `json:"frameworks"`
}

type Framwork struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	MarathonUrl string `json:"webui_url"`
}
