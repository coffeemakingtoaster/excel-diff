package fileutils

import (
	"errors"
	"os"
)

func FileExists(path string) bool {
	fi, err := os.Stat("path")
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	switch mode := fi.Mode(); {
	case mode.IsRegular():
		return true
	}
	return false
}
