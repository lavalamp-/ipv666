package data

import (
	"github.com/lavalamp-/ipv666/common/modeling"
	"log"
	"errors"
	"fmt"
	"path/filepath"
	"github.com/lavalamp-/ipv666/common/fs"
)

var curAddressModel modeling.ProbabilisticAddressModel
var curAddressModelPath string

func GetProbabilisticAddressModel(modelDir string) (modeling.ProbabilisticAddressModel, error) {
	log.Printf("Attempting to retrieve most recent probabilistic model from directory '%s'.", modelDir)
	fileName, err := fs.GetMostRecentFileFromDirectory(modelDir)
	if err != nil {
		log.Printf("Error thrown when retrieving probabilistic model from directory '%s': %s", modelDir, err)
		return modeling.ProbabilisticAddressModel{}, err
	} else if fileName == "" {
		log.Printf("The directory at '%s' was empty.", modelDir)
		return modeling.ProbabilisticAddressModel{}, errors.New(fmt.Sprintf("No model files were found in directory %s.", modelDir))
	}
	filePath := filepath.Join(modelDir, fileName)
	log.Printf("Most recent probabilistic address model is at path '%s'.", filePath)
	if fileName == curAddressModelPath {
		log.Printf("Already have model at path '%s' loaded in memory. Returning.", filePath)
		return curAddressModel, nil
	} else {
		log.Printf("Loading probabilistic address model from path '%s'.", filePath)
		return modeling.GetProbabilisticModelFromFile(filePath)
	}
}

func GetMostRecentCandidateFilePath(candidateDir string) (string, error) {
	log.Printf("Attempting to find most recent candidate file path in directory '%s'.", candidateDir)
	fileName, err := fs.GetMostRecentFileFromDirectory(candidateDir)
	if err != nil {
		log.Printf("Error thrown when finding most recent candidate file path in directory '%s': %s", candidateDir, err)
		return "", err
	} else if fileName == "" {
		return "", errors.New(fmt.Sprintf("No candidate file was found in directory '%s'.", candidateDir))
	} else {
		log.Printf("Most recent file path in directory '%s' is '%s'.", candidateDir, fileName)
		filePath := filepath.Join(candidateDir, fileName)
		return filePath, nil
	}
}

