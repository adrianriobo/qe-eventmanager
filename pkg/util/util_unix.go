// +build !windows

package util

import (
	"os"
)

func GetHomeDir() string {
	return os.Getenv("HOME")
}
