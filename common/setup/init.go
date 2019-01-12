package setup

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/fs"
	"github.com/lavalamp-/ipv666/common/statemachine"
	"log"
	"os"
)

func InitFilesystem() error {
	log.Print("Now initializing filesystem...")
	for _, dirPath := range config.GetAllDirectories() {
		err := fs.CreateDirectoryIfNotExist(dirPath)
		if err != nil {
			return err
		}
	}
	log.Printf("Initializing state file at '%s'.", config.GetStateFilePath())
	if _, err := os.Stat(config.GetStateFilePath()); os.IsNotExist(err) {
		log.Printf("State file does not exist at path '%s'. Creating now.", config.GetStateFilePath())
		err = statemachine.InitStateFile(config.GetStateFilePath())
		if err != nil {
			return err
		}
	} else {
		log.Printf("State file already exists at path '%s'.", config.GetStateFilePath())
	}
	log.Print("Local filesystem initialized.")
	return nil
}
