package util

import (
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"time"
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
