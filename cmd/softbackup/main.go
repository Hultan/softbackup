package main

import (
	"fmt"
	"os/exec"
	"path"
	"time"
)

const (
	constDateLayoutBackup = "20060102_0304"
	constConfigPath       = "/.config/softteam/softbackup/softbackup.config"
)

func main() {
	fmt.Println("Starting database backup!")
	fmt.Println("")

	// Load config file
	config := new(Config)
	err := config.Load()
	if err!=nil {
		fmt.Println("Failed to load ~/.config/softteam/softbackup/softbackup.config!")
		panic(err)
	}

	fmt.Println("Server:", config.Connection.Server)
	fmt.Println("User:", config.Connection.Username)
	fmt.Println("Databases to backup:", config.Databases)
	fmt.Println("Destination:", config.Destination.Path)
	fmt.Println("")

	for _, item := range config.Databases {
		fmt.Println("Starting backup up database :", item)
		err := backup(item, config.Destination.Path)
		if err != nil {
			fmt.Println("Error when backing up database :", item, ":", err.Error())
			return
		}
		
		fmt.Println("Finished backing up database :", item, "!")
	}
}

// Backs up a mysql database
func backup(database, rootBackupPath string) error {
	backupFile := fmt.Sprintf("%s_%s.sql", database, time.Now().Local().Format(constDateLayoutBackup))
	backupPath := path.Join(rootBackupPath, backupFile)
	command := fmt.Sprintf("mysqldump -u per %s > %s", database, backupPath)
	fmt.Println("Running command :", command)
	_, err := exec.Command("/bin/bash", "-c", command).Output()
	if err != nil {
		return err
	}

	return nil
}
