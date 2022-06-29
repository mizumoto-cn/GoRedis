package fileopt

import (
	"fmt"
	"os"
	"path/filepath"
)

// CheckPermission checks if the path has permission to open the given path (dir or file)
func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	// whether the error is known to report that permission is denied
	return os.IsPermission(err)
}

// CheckExistence checks if the path exists
func CheckExistence(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

// MkdirIfNotExist creates a directory if it does not exist
func MkdirIfNotExist(dir string) error {
	if CheckExistence(dir) {
		err := Mkdir(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

// Mkdir creates a directory
func Mkdir(dir string) error {
	err := os.MkdirAll(dir, 0777) // maybe 0755?
	if err != nil {
		return err
	}
	return nil
}

func SafeOpen(fileName, dir string) (*os.File, error) {
	permissionCheck := CheckPermission(dir)
	if permissionCheck {
		return nil, fmt.Errorf("permission denied: %w", os.ErrPermission)
	}

	if err := MkdirIfNotExist(dir); err != nil {
		return nil, fmt.Errorf("unable tp mkdir at %s : %w", dir, err)
	}

	file, err := os.OpenFile(filepath.Join(dir, fileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s : %w", fileName, err)
	}

	return file, nil
}
