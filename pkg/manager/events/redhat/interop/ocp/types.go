package ocp

import (
	"github.com/adrianriobo/qe-eventmanager/pkg/manager/events/redhat/interop"
)

type BuildComplete struct {
	Artifact   Artifact         `json:"artifact"`
	Contact    interop.Contact  `json:"contact"`
	GenerateAt string           `json:"generated_at"`
	System     []interop.System `json:"system"`
	Version    string           `json:"version"`
}

type TestComplete struct {
	Contact    interop.Contact  `json:"contact"`
	Run        interop.Run      `json:"run"`
	Artifact   Artifact         `json:"artifact"`
	Test       interop.Test     `json:"test"`
	GenerateAt string           `json:"generated_at"`
	System     []interop.System `json:"system"`
	Version    string           `json:"version"`
}

type TestError struct {
	Contact    interop.Contact  `json:"contact"`
	Run        interop.Run      `json:"run"`
	Artifact   Artifact         `json:"artifact"`
	Test       interop.Test     `json:"test"`
	Error      interop.Error    `json:"error"`
	GenerateAt string           `json:"generated_at"`
	System     []interop.System `json:"system"`
	Version    string           `json:"version"`
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
	Artifacts                []interface{} `json:"artifacts"`
	Phase                    string        `json:"phase"`
	Release                  string        `json:"release"`
}
