package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
)

// Config for the database-backup program
type Config struct {
	Connection struct {
		Server   string `json:"server"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"connection"`
	Destination struct {
		Path   string `json:"path"`
	} `json:"destination"`
	Databases []string `json:"databases"`
}

// Load : Loads a SoftTube configuration file
func (config *Config) Load() error {
	// Get the path to the config file
	path := getConfigPath()

	// Make sure the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		errorMessage := fmt.Sprintf("settings file is missing (%s)", constConfigPath)
		return errors.New(errorMessage)
	}

	// Open config file
	configFile, err := os.Open(path)

	// Handle errors
	if err != nil {
		fmt.Println(err.Error())
	}
	defer configFile.Close()

	// Parse the JSON document
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return nil
}

// Save : Saves a SoftTube configuration file
func (config *Config) Save(mode string) {
	// Get the path to the config file
	path := getConfigPath()

	// Open config file
	configFile, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, 0644)

	// Handle errors
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer configFile.Close()

	// Create JSON from config object
	data, err := json.MarshalIndent(config, "", "\t")

	// Handle errors
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Write the data
	configFile.Write(data)
}

// Get path to the config file
// Mode = "test" returns test config path
// otherwise returns normal config path
func getConfigPath() string {
	home := getHomeDirectory()
	configPath := constConfigPath

	return path.Join(home, configPath)
}

// Get current users home directory
func getHomeDirectory() string {
	u, err := user.Current()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get user home directory : %s", err)
		panic(errorMessage)
	}
	return u.HomeDir
}
