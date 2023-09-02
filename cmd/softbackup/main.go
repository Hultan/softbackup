package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/hultan/crypto"
	"github.com/hultan/logger"
)

const (
	applicationVersion    = "1.0.1"
	constDateLayoutBackup = "20060102_0304"
	constConfigPath       = "/.config/softteam/softbackup/softbackup.config"
	errorOpenConfig       = 1
	errorOpenLog          = 2
)

var (
	log    *logger.Logger
	config *Config
)

var servers map[string]Server

func main() {
	// Handle arguments
	if len(os.Args) > 1 {
		if strings.HasPrefix(os.Args[1], "-v") {
			fmt.Printf("SoftBackup - %s\n", applicationVersion)
		} else {
			fmt.Printf("Usage:\n")
			fmt.Printf("    softbackup [-v]|[-version]\n")
		}
		os.Exit(0)
	}

	// Load config file
	config = &Config{}
	err := config.Load()
	if err != nil {
		fmt.Printf("Failed to open config file! ('%s')", constConfigPath)
		fmt.Println(err)
		os.Exit(errorOpenConfig)
	}

	// Create server map
	servers = make(map[string]Server, len(config.Servers))
	for _, server := range config.Servers {
		servers[server.Name] = server
	}

	log, err = logger.NewStandardLogger(path.Join(config.Paths.Log, "softbackup.log"))
	if err != nil {
		fmt.Printf("Failed to open log file! ('%s')", config.Paths.Log)
		fmt.Println(err)
		os.Exit(errorOpenLog)
	}
	defer log.Close()

	// Start updating the softtube database
	log.Info.Println()
	log.Info.Println("-------------------")
	log.Info.Println("softtube.softbackup")
	log.Info.Println("-------------------")
	log.Info.Println()
	log.Info.Println("Configuration:")
	log.Info.Println()
	log.Info.Printf("Servers to backup: \n")
	for _, server := range config.Servers {
		log.Info.Printf("\t%s\n", server)
	}
	log.Info.Println()
	log.Info.Printf("Databases to backup: \n")
	for _, database := range config.Databases {
		log.Info.Printf("\t%s\n", database)
	}
	log.Info.Println()
	log.Info.Printf("Paths: \n")
	log.Info.Printf("\tBackup : %s\n", config.Paths.Backup)
	log.Info.Printf("\tLog    : %s\n", config.Paths.Log)
	log.Info.Println()

	for _, database := range config.Databases {
		log.Info.Printf("Starting backup up database : '%s'\n", database)
		output, err := backup(database, config.Paths.Backup)
		if err != nil {
			log.Error.Printf("Failed to back up database '%s' : %v\n", database, err)
			log.Error.Printf("Output : %s\n", output)
		} else {
			log.Info.Printf("Successfully backed up database '%s'...\n", database)
		}
		log.Info.Println()
	}

	log.Info.Println("Finished backing up databases!")
}

// Backs up a mysql database
func backup(database Database, rootBackupPath string) (string, error) {
	backupFile := fmt.Sprintf("%s_%s_%s.sql",
		database.Server,
		database.Database,
		time.Now().Local().Format(constDateLayoutBackup))

	backupPath := path.Join(rootBackupPath, backupFile)

	server := servers[database.Server]
	command, logCommand := getCommand(server, database, backupPath)

	// Make sure password is not exposed in log files
	log.Info.Printf("Executing command : /bin/bash -c %s\n", logCommand)
	cmd := exec.Command("/bin/bash", "-c", command)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return string(bytes), err
	}

	return "", nil
}

func getCommand(server Server, database Database, backupPath string) (command string, logCommand string) {
	if server.Password == "" {
		command = fmt.Sprintf(
			"mysqldump --host %s -P %s -u %s %s > %s",
			server.Address, server.Port, server.UserName, database.Database, backupPath,
		)
		logCommand = command
	} else {
		password, err := crypto.Decrypt(server.Password)
		if err != nil {
			panic(err)
		}
		command = fmt.Sprintf(
			"mysqldump --host %s -P %s -u %s -p%s %s > %s",
			server.Address, server.Port, server.UserName, password, database.Database, backupPath,
		)
		logCommand = fmt.Sprintf(
			"mysqldump --host %s -P %s -u %s -p****** %s > %s",
			server.Address, server.Port, server.UserName, database.Database, backupPath,
		)
	}
	return
}
