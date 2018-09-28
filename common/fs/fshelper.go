package fs

import (
	"log"
	"os"
	"io/ioutil"
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

func GetMostRecentFileFromDirectory(dirPath string) (string, error) {

	// https://stackoverflow.com/questions/45578172/golang-find-most-recent-file-by-date-and-time

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Printf("Error thrown when reading files from directory '%s': %s", dirPath, err)
		return "", err
	}
	var newestFile string = ""
	var newestTime int64 = 0
	for _, fi := range files {
		if fi.Mode().IsRegular() {
			curTime := fi.ModTime().Unix()
			if curTime > newestTime {
				newestTime = curTime
				newestFile = fi.Name()
			}
		}
	}
	return newestFile, nil
}
