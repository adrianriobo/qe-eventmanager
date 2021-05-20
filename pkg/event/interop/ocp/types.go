package ocp

type BuildComplete struct {
	Artifact   Artifact `json:"artififact"`
	Contact    Contact  `json:"contact"`
	GenerateAt string   `json:"generated_at"`
	System     []System `json:"system"`
	Version    string   `json:"version"`
}

type TestComplete struct {
	Contact    Contact  `json:"contact"`
	Run        Run      `json:"run"`
	Artifact   Artifact `json:"artififact"`
	Test       Test     `json:"test"`
	GenerateAt string   `json:"generated_at"`
	System     []System `json:"system"`
	Version    string   `json:"version"`
}

type TestError struct {
	Contact    Contact  `json:"contact"`
	Run        Run      `json:"run"`
	Artifact   Artifact `json:"artififact"`
	Test       Test     `json:"test"`
	Error      Error    `json:"error"`
	GenerateAt string   `json:"generated_at"`
	System     []System `json:"system"`
	Version    string   `json:"version"`
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

type Run struct {
	URL string `json:"url"`
	Log string `json:"log"`
}

type Test struct {
	Category  string   `json:"category"`
	Namespace string   `json:"namespace"`
	TestType  string   `json:"type"`
	Result    string   `json:"result"`
	Runtime   string   `json:"runtime"`
	XunitUrls []string `json:"xunit_urls"`
}

type Error struct {
	Reason string `json:"reason"`
}
