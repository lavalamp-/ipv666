package data

import (
	"github.com/lavalamp-/ipv666/common/modeling"
	"log"
	"errors"
	"fmt"
	"path/filepath"
	"github.com/lavalamp-/ipv666/common/fs"
	"github.com/lavalamp-/ipv666/common/addressing"
	"net"
)

var curAddressModel modeling.ProbabilisticAddressModel
var curAddressModelPath string
var curCandidatePingResults addressing.IPv6AddressList
var curCandidatePingResultsPath string
var curScanResultsNetworkRanges []*net.IPNet
var curScanResultsNetworkRangesPath string

func UpdateScanResultsNetworkRanges(networks []*net.IPNet, filePath string) {
	curScanResultsNetworkRanges = networks
	curScanResultsNetworkRangesPath = filePath
}

func GetScanResultsNetworkRanges(scanResultsDir string) ([]*net.IPNet, error) {
	log.Printf("Attempting to retrieve most recent candidate ping networks from directory '%s'.", scanResultsDir)
	fileName, err := fs.GetMostRecentFileFromDirectory(scanResultsDir)
	if err != nil {
		log.Printf("Error thrown when retrieving candidate ping networks from directory '%s': %s", scanResultsDir, err)
		return nil, err
	} else if fileName == "" {
		log.Printf("The directory at '%s' was empty.", scanResultsDir)
		return nil, errors.New(fmt.Sprintf("No candidate ping networks files were found in directory %s.", scanResultsDir))
	}
	filePath := filepath.Join(scanResultsDir, fileName)
	log.Printf("Most recent candidate ping networks file is at path '%s'.", filePath)
	if fileName == curScanResultsNetworkRangesPath {
		log.Printf("Already have candidate ping networks at path '%s' loaded in memory. Returning.", filePath)
		return curScanResultsNetworkRanges, nil
	} else {
		log.Printf("Loading candidate ping networks from path '%s'.", filePath)
		toReturn, err := addressing.ReadIPv6NetworksFromFile(filePath)
		if err != nil {
			UpdateScanResultsNetworkRanges(toReturn, filePath)
		}
		return toReturn, err
	}
}

func GetCandidatePingResults(pingResultsDir string) (addressing.IPv6AddressList, error) {
	log.Printf("Attempting to retrieve most recent candidate ping results from directory '%s'.", pingResultsDir)
	fileName, err := fs.GetMostRecentFileFromDirectory(pingResultsDir)
	if err != nil {
		log.Printf("Error thrown when retrieving candidate ping results from directory '%s': %s", pingResultsDir, err)
		return addressing.IPv6AddressList{}, err
	} else if fileName == "" {
		log.Printf("The directory at '%s' was empty.", pingResultsDir)
		return addressing.IPv6AddressList{}, errors.New(fmt.Sprintf("No candidate ping files were found in directory %s.", pingResultsDir))
	}
	filePath := filepath.Join(pingResultsDir, fileName)
	log.Printf("Most recent ping results file is at path '%s'.", filePath)
	if fileName == curCandidatePingResultsPath {
		log.Printf("Already have candidate ping results at path '%s' loaded in memory. Returning.", filePath)
		return curCandidatePingResults, nil
	} else {
		log.Printf("Loading candidate ping results from path '%s'.", filePath)
		toReturn, err := addressing.GetAddressListFromHexStringsFile(filePath)
		if err != nil {
			curCandidatePingResultsPath = filePath
			curCandidatePingResults = toReturn
		}
		return toReturn, err
	}
}

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
		toReturn, err := modeling.GetProbabilisticModelFromFile(filePath)
		if err != nil {
			curAddressModelPath = filePath
			curAddressModel = toReturn
		}
		return toReturn, err
	}
}

func GetMostRecentFilePathFromDir(candidateDir string) (string, error) {
	log.Printf("Attempting to find most recent file path in directory '%s'.", candidateDir)
	fileName, err := fs.GetMostRecentFileFromDirectory(candidateDir)
	if err != nil {
		log.Printf("Error thrown when finding most recent candidate file path in directory '%s': %s", candidateDir, err)
		return "", err
	} else if fileName == "" {
		return "", errors.New(fmt.Sprintf("No file was found in directory '%s'.", candidateDir))
	} else {
		log.Printf("Most recent file path in directory '%s' is '%s'.", candidateDir, fileName)
		filePath := filepath.Join(candidateDir, fileName)
		return filePath, nil
	}
}
