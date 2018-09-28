package common

import (
	"log"
	"os"
)

func CreateDirectoryIfNotExist(dirPath string) (error) {
	log.Printf("Making sure that directory at '%s' exists.", dirPath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		log.Printf("No directory found at path '%s'. Creating now.", dirPath)
		return os.Mkdir(dirPath, 0755)
	} else {
		log.Printf("Directory at path '%s' already exists.", dirPath)
		return nil
	}
}
