// Package config provides functions for caching state
package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// The base name of mz's configuration file
const configFileName = "config.toml"

// Saved program state
type last struct {
	Name string
}

// Configuration allows for saving program configuraiton and settings
type Configuration struct {
	Last      last
	Favorites []string
}

// Save writes the configuration to disk
func (s *Configuration) Save() error {
	file, err := GetSettingsFile()
	if err != nil {
		return err
	}

	dat, err := toml.Marshal(s)
	if err != nil {
		return err
	}

	err = os.WriteFile(file, dat, 0644)
	if err != nil {
		return err
	}

	return nil
}

// GetSettingsFile returns the path of mz's cache file.
// The file and its parent directories may not exist.
func GetSettingsFile() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil && !os.IsExist(err) {
		return "", err
	}

	folder := filepath.Join(configDir, "mz")

	err = os.MkdirAll(folder, 0755)
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(folder, configFileName)

	return configPath, nil
}

// LoadConfig attempts to load the existing configuration file.
// If an error occurs (i.e. the configuration file does not exist)
// an empty configuration struct is returned.
func LoadConfig() *Configuration {
	var s Configuration

	file, err := GetSettingsFile()
	if err != nil {
		log.Println(err)
		return &s
	}

	dat, err := os.ReadFile(file)
	if err != nil {
		log.Println(err)
		return &s
	}

	err = toml.Unmarshal(dat, &s)
	if err != nil {
		log.Println(err)
		return &s
	}

	return &s
}
