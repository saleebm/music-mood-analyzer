package shared

import (
	"errors"
	"os"
)

func Exists(name string) (bool, error) {
	info, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return !info.IsDir(), err
}
