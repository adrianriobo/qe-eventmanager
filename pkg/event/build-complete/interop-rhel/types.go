package interopRHEL

import (
	buildComplete "github.com/adrianriobo/qe-eventmanager/pkg/event/build-complete"
)

type BuildComplete struct {
	Artifact   Artifact               `json:"artifact"`
	Contact    buildComplete.Contact  `json:"contact"`
	GenerateAt string                 `json:"generated_at"`
	System     []buildComplete.System `json:"system"`
	Version    string                 `json:"version"`
}

type TestComplete struct {
	Contact    buildComplete.Contact  `json:"contact"`
	Run        buildComplete.Run      `json:"run"`
	Artifact   Artifact               `json:"artifact"`
	Test       buildComplete.Test     `json:"test"`
	GenerateAt string                 `json:"generated_at"`
	System     []buildComplete.System `json:"system"`
	Version    string                 `json:"version"`
}

type TestError struct {
	Contact    buildComplete.Contact  `json:"contact"`
	Run        buildComplete.Run      `json:"run"`
	Artifact   Artifact               `json:"artifact"`
	Test       buildComplete.Test     `json:"test"`
	Error      buildComplete.Error    `json:"error"`
	GenerateAt string                 `json:"generated_at"`
	System     []buildComplete.System `json:"system"`
	Version    string                 `json:"version"`
}

type Artifact struct {
	ArtifcatType string    `json:"type"`
	Id           string    `json:"id"`
	Products     []Product `json:"products"`
	Email        string    `json:"email"`
	Url          string    `json:"url"`
}

type Product struct {
	Architecture string        `json:"architecture"`
	Artifacts    []interface{} `json:"artifacts"`
	Build        string        `json:"build"`
	Id           string        `json:"id"`
	Image        string        `json:"image"`
	Name         string        `json:"name"`
	NVR          string        `json:"nvr"`
	Phase        string        `json:"phase"`
	Release      string        `json:"release"`
	Repos        []Repository  `json:"repos"`
	State        string        `json:"state"`
	ProductType  string        `json:"type"`
	Version      string        `json:"version"`
}

type Repository struct {
	Base_Url string `json:"base_url"`
	Name     string `json:"name"`
}
