// Package config provides functions for caching state
package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

// GetCacheFile returns the path of mz's cache file.
// The file and its parent directories may not exist.
func GetCacheFile() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	cacheFile := filepath.Join(cacheDir, "mz")

	return cacheFile, nil
}

// SaveLastCanteen saves the name of the last viewed cantine.
func SaveLastCanteen(name string) error {
	file, err := GetCacheFile()
	if err != nil {
		return err
	}

	dat := []byte(name)

	err = os.WriteFile(file, dat, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// GetLastCanteen returns the name of the last viewed canteen.
// If an error occurs or the mensa is not found, an empty string will be returned.
func GetLastCanteen() string {
	file, err := GetCacheFile()
	if err != nil {
		log.Print(err)
		return ""
	}

	dat, err := os.ReadFile(file)
	if err != nil {
		// Ignore missing files
		if !errors.Is(err, os.ErrNotExist) {
			log.Print(err)
		}
		return ""
	}

	return string(dat)
}
