package util

import (
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"time"

	"golang.org/x/exp/slices"
)

func EnsureBaseDirectoriesExist(path string) error {
	return os.MkdirAll(path, 0750)
}

func GetHomeDir() string {
	return os.Getenv("HOME")
}

func GenerateCorrelation() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(1000000 + rand.Intn(8000000))
}

func IsEmpty(source interface{}) bool {
	nonEmptyValues := []interface{}(nil)
	structIterator := reflect.ValueOf(source)
	for i := 0; i < structIterator.NumField(); i++ {
		val := structIterator.Field(i).Interface()
		if !reflect.DeepEqual(val, reflect.Zero(structIterator.Field(i).Type()).Interface()) {
			nonEmptyValues = append(nonEmptyValues, val)
		}
	}
	return len(nonEmptyValues) == 0
}

func SliceContains[T comparable](sourceSlice []T, item T) bool {
	idx := slices.IndexFunc(sourceSlice,
		func(e T) bool { return e == item })
	return idx != -1
}

// For an slice of any Type, we can check if is there any element based on
// a custom contains function
// if contains is true we can get the element based on a custom value function
func SliceItem[T any](sourceSlice []T,
	contains func(e T) bool,
	value func(e T) T) (bool, T) {
	idx := slices.IndexFunc(sourceSlice,
		func(e T) bool { return contains(e) })
	if idx == -1 {
		return false, getZero[T]()
	}
	return true, value(sourceSlice[idx])
}

func getZero[T any]() T {
	var result T
	return result
}
