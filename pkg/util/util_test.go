package util

import (
	"testing"
)

func TestSliceContains(t *testing.T) {
	contains := SliceContains([]string{"foo", "bar"}, "foo")
	if !contains {
		t.Fatal("Slice shoud contain")
	}
}

func TestSliceNotContains(t *testing.T) {
	contains := SliceContains([]string{"foo", "bar"}, "baz")
	if contains {
		t.Fatal("Slice shoud not contain")
	}
}

func TestSliceItem(t *testing.T) {
	contains, value := SliceItem([]string{"foo", "bar"}, "foo", func(e string) bool { return e == "foo" }, func(e string) string { return e })
	if !contains || value != "foo" {
		t.Fatal("Slice shoud contain")
	}
}
