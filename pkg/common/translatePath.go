package common

import (
	"path/filepath"
)

func AbsPath(path *string) error {
	if len(*path) <= 0 {
		return nil
	}

	var err error
	*path, err = filepath.Abs(*path)
	return err
}
