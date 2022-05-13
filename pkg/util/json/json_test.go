package json

import (
	"testing"
)

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

func TestMatching(t *testing.T) {
	match, err := MatchFilters([]byte(event), []string{"$.artifact.products[?(@.nvr=='found')].nvr"})
	if !match || err != nil {
		t.Fatal("Expression should match")
	}
}

func TestNotMatching(t *testing.T) {

	match, err := MatchFilters([]byte(event), []string{"$.artifact.products[?(@.nvr=='not-found')].nvr"})
	if match || err == nil {
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
