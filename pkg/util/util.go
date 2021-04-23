package util

import "os"

func EnsureBaseDirectoriesExist(path string) error {
	return os.MkdirAll(path, 0750)
}
