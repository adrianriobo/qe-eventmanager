// +build windows

package util

import (
	"os"
	"path/filepath"
)

func GetHomeDir() string {
	if homeDrive, homePath := os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"); len(homeDrive) > 0 && len(homePath) > 0 {
		homeDir := filepath.Join(homeDrive, homePath)
		if _, err := os.Stat(homeDir); err == nil {
			return homeDir
		}
	}
	if userProfile := os.Getenv("USERPROFILE"); len(userProfile) > 0 {
		if _, err := os.Stat(userProfile); err == nil {
			return userProfile
		}
	}
	return os.Getenv("HOME")
}
