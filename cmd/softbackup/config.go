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
	Paths struct {
		Backup string `json:"backup"`
		Log    string `json:"log"`
	} `json:"paths"`
	Servers   []Server   `json:"servers"`
	Databases []Database `json:"databases"`
}

type Database struct {
	Server   string `json:"server"`
	Database string `json:"database"`
}

type Server struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Port     string `json:"port"`
}

func (s Server) String() string {
	return fmt.Sprintf("%-10s: %-13s (port:%s, user:%s)", s.Name, s.Address, s.Port, s.UserName)
}

func (d Database) String() string {
	return fmt.Sprintf("%-10s: %s", d.Server, d.Database)
}

// Load : Loads a SoftBackup configuration file
func (config *Config) Load() error {
	// Get the path to the config file
	configPath := getConfigPath()

	// Make sure the file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		errorMessage := fmt.Sprintf("settings file is missing (%s)", constConfigPath)
		return errors.New(errorMessage)
	}

	// Open config file
	configFile, err := os.Open(configPath)

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

// Get path to the config file
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
