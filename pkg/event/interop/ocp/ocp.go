package ocp

type Event struct {
	Message Message `json:"message"`
	Topic   string  `json:"topic"`
}

type Message struct {
	Contact    Contact  `json:"contact"`
	GenerateAt string   `json:"generated_at"`
	Version    string   `json:"version"`
	Artifact   Artifact `json:"artififact"`
	System     []System `json:"system"`
}

type Contact struct {
	Name  string `json:"name"`
	Team  string `json:"team"`
	Docs  string `json:"docs"`
	Email string `json:"email"`
	Url   string `json:"url"`
}

type Artifact struct {
	ArtifcatType string    `json:"type"`
	Id           string    `json:"id"`
	Products     []Product `json:"products"`
	Email        string    `json:"email"`
	Url          string    `json:"url"`
}

type Product struct {
	Id                       string        `json:"id"`
	NVR                      string        `json:"nvr"`
	Name                     string        `json:"name"`
	Version                  string        `json:"version"`
	Architecture             string        `json:"architecture"`
	Build                    string        `json:"build"`
	Internal_build_index_url string        `json:"internal_build_index_url"`
	External_build_index_url string        `json:"external_build_index_url"`
	ProductType              string        `json:"type"`
	State                    string        `json:"state"`
	Artifacts                []interface{} `json:"artifcats"`
	Phase                    string        `json:"phase"`
	Release                  string        `json:"release"`
}

type System struct {
	Architecture string `json:"architecture"`
	Provider     string `json:"provider"`
	OS           string `json:"os"`
}
