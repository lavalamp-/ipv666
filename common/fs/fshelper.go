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

func GetNonMostRecentFilesFromDirectory(dirPath string) ([]string, error) {
	var toReturn []string
	recentFile, err := GetMostRecentFileFromDirectory(dirPath)
	if err != nil || recentFile == ""{
		return toReturn, err
	}
	log.Printf("Most recent file in directory '%s' is '%s'.", dirPath, recentFile)
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Printf("Error thrown when trying to read files from directory '%s': '%s", dirPath, err)
		return toReturn, err
	}
	for _, fi := range files {
		name := fi.Name()
		if name != recentFile {
			toReturn = append(toReturn, name)
		}
	}
	log.Printf("Found %d files older than the most recent '%s' in directory '%s'.", len(toReturn), recentFile, dirPath)
	return toReturn, nil
}
