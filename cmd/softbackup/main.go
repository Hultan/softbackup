package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/hultan/softteam/framework"
)

const (
	constDateLayoutBackup = "20060102_0304"
	constConfigPath       = "/.config/softteam/softbackup/softbackup.config"
	errorOpenConfig       = 1
	errorOpenLog          = 2
)

var (
	fw     *framework.Framework
	logger *framework.Logger
	config *Config
)

func main() {
	// Load config file
	config = &Config{}
	err := config.Load()
	if err != nil {
		fmt.Printf("Failed to open config file! ('%s')", constConfigPath)
		fmt.Println(err)
		os.Exit(errorOpenConfig)
	}

	fw = framework.NewFramework()
	logger, err = fw.Log.NewStandardLogger(path.Join(config.Paths.Log, "softbackup.log"))
	if err != nil {
		fmt.Printf("Failed to open log file! ('%s')", config.Paths.Log)
		fmt.Println(err)
		os.Exit(errorOpenLog)
	}
	defer logger.Close()

	// Start updating the softtube database
	logger.Info.Println()
	logger.Info.Println("-------------------")
	logger.Info.Println("softtube.softbackup")
	logger.Info.Println("-------------------")
	logger.Info.Println()
	logger.Info.Println("Configuration:")
	logger.Info.Println()
	logger.Info.Printf("Server             : %s\n", config.Connection.Server)
	logger.Info.Printf("User               : %s\n", config.Connection.Username)
	logger.Info.Printf("Databases to backup: %s\n", config.Databases)
	logger.Info.Printf("Backup destination : %s\n", config.Paths.Backup)
	logger.Info.Printf("Log destination    : %s\n", config.Paths.Log)
	logger.Info.Println()

	for _, database := range config.Databases {
		logger.Info.Printf("Starting backup up database : '%s'\n", database)
		output, err := backup(database, config.Paths.Backup)
		if err != nil {
			logger.Error.Printf("Failed to back up database '%s' : %v\n", database, err)
			logger.Error.Printf("Output : %s\n", output)
		} else {
			logger.Info.Printf("Successfully backed up database '%s'...\n", database)
		}
		logger.Info.Println()
	}

	logger.Info.Println("Finished backing up databases!")
}

// Backs up a mysql database
func backup(database, rootBackupPath string) (string, error) {
	backupFile := fmt.Sprintf("%s_%s.sql", database, time.Now().Local().Format(constDateLayoutBackup))
	backupPath := path.Join(rootBackupPath, backupFile)
	command := fmt.Sprintf("mysqldump -u per %s > %s", database, backupPath)
	logger.Info.Printf("Executing command : /bin/bash -c %s\n", command)
	bytes, err := exec.Command("/bin/bash", "-c", command).Output()
	if err != nil {
		return string(bytes), err
	}

	return "", nil
}
