package json

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
)

type Artifact struct {
	Products []Product `json:"products"`
}

type Product struct {
	ID  string `json:"id"`
	NVR string `json:"nvr"`
}

const event = `
{
	"contact": {
		"name": "foo"
	},
	"artifact": {
		"products": [{
			"id": "",
			"nvr": "found"
		}, {
			"id": "foo",
			"nvr": "bar"
		}]
	},
	"system": [{
		"architecture": "x86_64"
	}]
}`

const simpleEvent = `
{
	"action": "synchronize",
	"number": 3179
}`

func TestMatching(t *testing.T) {
	match, err := MatchFilters([]byte(event), []string{"$.artifact.products[?(@.nvr=='found')].nvr"})
	if !match || err != nil {
		t.Fatal("Expression should match")
	}
}

func TestMatchingSimple(t *testing.T) {
	match, err := MatchFilters([]byte(simpleEvent), []string{"$[?($.action == 'synchronize' || $.action == 'opened')]"})
	if !match || err != nil {
		t.Fatal("Expression should match")
	}
}

func TestNotMatching(t *testing.T) {
	match, err := MatchFilters([]byte(event), []string{"$.artifact.products[?(@.nvr=='not-found')].nvr"})
	if match || err != nil {
		t.Fatal("Expression should not match")
	}
}

func TestGetStringValueFound(t *testing.T) {
	value, err := GetStringValue([]byte(event), "$.artifact.products[?(@.nvr=='bar')].id")
	expectedValue := "foo"
	if value != expectedValue || err != nil {
		t.Fatal("Expression should match")
	}
}

func TestGetStringValueNotFound(t *testing.T) {
	value, err := GetStringValue([]byte(event), "$.artifact.products[?(@.nvr=='not-found')].id")
	expectedValue := ""
	if value != expectedValue || err == nil {
		t.Fatal("Expression should not match")
	}
}

func TestGetNodeValueFound(t *testing.T) {
	value, err := GetNodeAsByteArray([]byte(event), "$.artifact")
	valueAsString := trim(value)
	logging.Debugf("value for node %s", valueAsString)
	expectedValue := `{
	"products": [{
		"id": "",
		"nvr": "found"
	}, {
		"id": "foo",
		"nvr": "bar"
	}]}`
	planeExpected := trim([]byte(expectedValue))
	if valueAsString != planeExpected || err != nil {
		t.Fatal("Expression should match")
	}
}

func trim(event []byte) string {
	eventWithoutN := strings.ReplaceAll(string(event), "\n", "")
	return strings.ReplaceAll(eventWithoutN, "\t", "")
}

func TestMarshalling(t *testing.T) {
	value, err := GetNodeAsByteArray([]byte(event), "$.artifact")
	var artifact Artifact
	if err == nil {
		logging.Debugf("value for node %s", string(value))
		err = json.Unmarshal(value, &artifact)
	}
	if len(artifact.Products) != 2 || err != nil {
		t.Fatal("Expression should match")
	}
}
