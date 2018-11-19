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
	"github.com/lavalamp-/ipv666/common/blacklist"
	"github.com/willf/bloom"
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/filtering"
	"os"
	"io/ioutil"
)

var curAddressModel *modeling.ProbabilisticAddressModel
var curAddressModelPath string
var curCandidatePingResults []*net.IP
var curCandidatePingResultsPath string
var curScanResultsNetworkRanges []*net.IPNet
var curScanResultsNetworkRangesPath string
var curBlacklist *blacklist.NetworkBlacklist
var curBlacklistPath string
var curCleanPingResults []*net.IP
var curCleanPingResultsPath string
var curBloomFilter *bloom.BloomFilter
var curBloomFilterPath string
var curAliasedNetworks []*net.IPNet
var curAliasedNetworksPath string


func GetMostRecentTargetNetworkString(conf *config.Configuration) (string, error) {
	if !fs.CheckIfFileExists(conf.GetTargetNetworkFilePath()) {
		return "", nil
	}
	content, err := ioutil.ReadFile(conf.GetTargetNetworkFilePath())
	if err != nil {
		return "", nil
	}
	network, err := addressing.GetIPv6NetworkFromBytesIncLength(content)
	if err != nil {
		return "", err
	}
	return network.String(), nil
}

func WriteMostRecentTargetNetwork(toWrite *net.IPNet, conf *config.Configuration) (error) {
	return addressing.WriteIPv6NetworksToFile(conf.GetTargetNetworkFilePath(), []*net.IPNet{toWrite})
}

func UpdateAliasedNetworks(nets []*net.IPNet, filePath string) {
	curAliasedNetworks = nets
	curAliasedNetworksPath = filePath
}

func GetAliasedNetworks(conf *config.Configuration) ([]*net.IPNet, error) {
	aliasedDir := conf.GetAliasedNetworkDirPath()
	log.Printf("Attempting to retrieve most recent aliased networks from directory '%s'.", aliasedDir)
	fileName, err := fs.GetMostRecentFileFromDirectory(aliasedDir)
	if err != nil {
		log.Printf("Error thrown when retrieving aliased networks from directory '%s': %s", aliasedDir, err)
		return nil, err
	} else if fileName == "" {
		log.Printf("The directory at '%s' was empty.", aliasedDir)
		return nil, errors.New(fmt.Sprintf("No aliased networks were found in directory %s.", aliasedDir))
	}
	filePath := filepath.Join(aliasedDir, fileName)
	log.Printf("Most recent aliased networks file is at path '%s'.", filePath)
	if filePath == curAliasedNetworksPath {
		log.Printf("Already have aliased networks from path '%s' loaded in memory. Returning.", filePath)
		return curAliasedNetworks, nil
	} else {
		log.Printf("Loading aliased networks from path '%s'.", filePath)
		toReturn, err := addressing.ReadIPv6NetworksFromFile(filePath)
		if err != nil {
			UpdateAliasedNetworks(toReturn, filePath)
		}
		return toReturn, err
	}
}

func UpdateBloomFilter(filter *bloom.BloomFilter, filePath string) {
	curBloomFilter = filter
	curBloomFilterPath = filePath
}

func LoadBloomFilterFromOutput(conf *config.Configuration) (*bloom.BloomFilter, error) {
	log.Printf("Creating Bloom filter from output file '%s'.", conf.GetOutputFilePath())
	ips, err := addressing.ReadIPsFromHexFile(conf.GetOutputFilePath())
	ips = addressing.GetUniqueIPs(ips, conf.LogLoopEmitFreq)
	if err != nil {
		return nil, err
	}
	log.Printf("%d IP addresses loaded from file '%s'.", len(ips), conf.GetOutputFilePath())
	newBloom := bloom.New(conf.AddressFilterSize, conf.AddressFilterHashCount)
	for _, ip := range ips {
		ipBytes := ([]byte)(*ip)
		newBloom.Add(ipBytes)
	}
	log.Printf("Created Bloom filter with %d addresses from '%s'.", len(ips), conf.GetOutputFilePath())
	return newBloom, nil
}

func GetBloomFilter(conf *config.Configuration) (*bloom.BloomFilter, error) {
	filterDir := conf.GetBloomDirPath()
	log.Printf("Attempting to retrieve most recent Bloom filter from directory '%s'.", filterDir)
	fileName, err := fs.GetMostRecentFileFromDirectory(filterDir)
	if err != nil {
		log.Printf("Error thrown when retrieving Bloom filter from directory '%s': %s", filterDir, err)
		return nil, err
	} else if fileName == "" {
		log.Printf("The directory at '%s' was empty. Checking for pre-existing output file at '%s'.", filterDir, conf.GetOutputFilePath())
		if _, err := os.Stat(conf.GetOutputFilePath()); !os.IsNotExist(err) {
			log.Printf("File at path '%s' exists. Using for new Bloom filter.", conf.GetOutputFilePath())
			return LoadBloomFilterFromOutput(conf)
		} else {
			log.Printf("No existing output file at '%s'. Returning a new, empty Bloom filter.", conf.GetOutputFilePath())
			return bloom.New(conf.AddressFilterSize, conf.AddressFilterHashCount), nil
		}
	}
	filePath := filepath.Join(filterDir, fileName)
	log.Printf("Most recent Bloom filter is at path '%s'.", filePath)
	if filePath == curBloomFilterPath {
		log.Printf("Already have Bloom filter at path '%s' loaded in memory. Returning.", filePath)
		return curBloomFilter, nil
	} else {
		log.Printf("Loading Bloom filter from path '%s'.", filePath)
		toReturn, err := filtering.GetBloomFilterFromFile(filePath, conf.AddressFilterSize, conf.AddressFilterHashCount)
		if err != nil {
			UpdateBloomFilter(toReturn, filePath)
		}
		return toReturn, err
	}
}

func UpdateCleanPingResults(addrs []*net.IP, filePath string) {
	curCleanPingResults = addrs
	curCleanPingResultsPath = filePath
}

func GetCleanPingResults(resultsDir string) ([]*net.IP, error) {
	log.Printf("Attempting to retrieve most recent cleaned ping results from directory '%s'.", resultsDir)
	fileName, err := fs.GetMostRecentFileFromDirectory(resultsDir)
	if err != nil {
		log.Printf("Error thrown when retrieving cleaned ping results from directory '%s': %e", resultsDir, err)
		return nil, err
	} else if fileName == "" {
		log.Printf("The directory at '%s' was empty.", resultsDir)
		return nil, errors.New(fmt.Sprintf("No cleaned ping results files were found in directory %s.", resultsDir))
	}
	filePath := filepath.Join(resultsDir, fileName)
	log.Printf("Most recent cleaned ping results file is at path '%s'.", filePath)
	if filePath == curCleanPingResultsPath {
		log.Printf("Already have cleaned ping results at path '%s' loaded in memory. Returning.", filePath)
		return curCleanPingResults, nil
	} else {
		log.Printf("Loading cleaned ping results from path '%s'.", filePath)
		toReturn, err := addressing.ReadIPsFromBinaryFile(filePath)
		if err != nil {
			UpdateCleanPingResults(toReturn, filePath)
		}
		return toReturn, err
	}
}

func UpdateBlacklist(blacklist *blacklist.NetworkBlacklist, filePath string) {
	curBlacklist = blacklist
	curBlacklistPath = filePath
}

func GetBlacklist(blacklistDir string) (*blacklist.NetworkBlacklist, error) {
	log.Printf("Attempting to retrieve most recent blacklist from directory '%s'.", blacklistDir)
	fileName, err := fs.GetMostRecentFileFromDirectory(blacklistDir)
	if err != nil {
		log.Printf("Error thrown when retrieving blacklist from directory '%s': %s", blacklistDir, err)
		return nil, err
	} else if fileName == "" {
		log.Printf("The directory at '%s' was empty. Returning a new, empty blacklist.", blacklistDir)
		emptyNets := make([]*net.IPNet, 0)
		return blacklist.NewNetworkBlacklist(emptyNets), nil
	}
	filePath := filepath.Join(blacklistDir, fileName)
	log.Printf("Most recent blacklist file is at path '%s'.", filePath)
	if filePath == curBlacklistPath {
		log.Printf("Already have blacklist at path '%s' loaded in memory. Returning.", filePath)
		return curBlacklist, nil
	} else {
		toReturn, err := blacklist.ReadNetworkBlacklistFromFile(filePath)
		if err != nil {
			UpdateBlacklist(toReturn, filePath)
		}
		return toReturn, err
	}
}

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
	if filePath == curScanResultsNetworkRangesPath {
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

func UpdateCandidatePingResults(ips []*net.IP, filePath string) {
	curCandidatePingResultsPath = filePath
	curCandidatePingResults = ips
}

func GetCandidatePingResults(pingResultsDir string) ([]*net.IP, error) {
	log.Printf("Attempting to retrieve most recent candidate ping results from directory '%s'.", pingResultsDir)
	fileName, err := fs.GetMostRecentFileFromDirectory(pingResultsDir)
	if err != nil {
		log.Printf("Error thrown when retrieving candidate ping results from directory '%s': %s", pingResultsDir, err)
		return nil, err
	} else if fileName == "" {
		log.Printf("The directory at '%s' was empty.", pingResultsDir)
		return nil, errors.New(fmt.Sprintf("No candidate ping files were found in directory %s.", pingResultsDir))
	}
	filePath := filepath.Join(pingResultsDir, fileName)
	log.Printf("Most recent ping results file is at path '%s'.", filePath)
	if filePath == curCandidatePingResultsPath {
		log.Printf("Already have candidate ping results at path '%s' loaded in memory. Returning.", filePath)
		return curCandidatePingResults, nil
	} else {
		log.Printf("Loading candidate ping results from path '%s'.", filePath)
		toReturn, err := addressing.ReadIPsFromHexFile(filePath)
		if err != nil {
			UpdateCandidatePingResults(toReturn, filePath)
		}
		return toReturn, err
	}
}

func UpdateProbabilisticAddressModel(model *modeling.ProbabilisticAddressModel, filePath string) {
	curAddressModelPath = filePath
	curAddressModel = model
}

func GetProbabilisticAddressModel(modelDir string) (*modeling.ProbabilisticAddressModel, error) {
	log.Printf("Attempting to retrieve most recent probabilistic model from directory '%s'.", modelDir)
	fileName, err := fs.GetMostRecentFileFromDirectory(modelDir)
	if err != nil {
		log.Printf("Error thrown when retrieving probabilistic model from directory '%s': %s", modelDir, err)
		return &modeling.ProbabilisticAddressModel{}, err
	} else if fileName == "" {
		log.Printf("The directory at '%s' was empty.", modelDir)
		return &modeling.ProbabilisticAddressModel{}, errors.New(fmt.Sprintf("No model files were found in directory %s.", modelDir))
	}
	filePath := filepath.Join(modelDir, fileName)
	log.Printf("Most recent probabilistic address model is at path '%s'.", filePath)
	if filePath == curAddressModelPath {
		log.Printf("Already have model at path '%s' loaded in memory. Returning.", filePath)
		return curAddressModel, nil
	} else {
		log.Printf("Loading probabilistic address model from path '%s'.", filePath)
		toReturn, err := modeling.GetProbabilisticModelFromFile(filePath)
		if err != nil {
			UpdateProbabilisticAddressModel(toReturn, filePath)
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
