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
			"id": "bar",
			"nvr": "bar"
		}]
	},
	"system": [{
		"architecture": "x86_64"
	}]
}`

func TestMatching(t *testing.T) {
	match, err := matchFilters(event, []string{"$.artifact.products[?(@.nvr=='found')].nvr"})
	if !match || err != nil {
		t.Fatal("Expression should match")
	}
}

func TestNotMatching(t *testing.T) {
	match, err := matchFilters(event, []string{"$.artifact.products[?(@.nvr=='not-found')].nvr"})
	if match || err != nil {
		t.Fatal("Expression should not match")
	}
}
