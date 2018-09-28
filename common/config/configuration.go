package config

import (
	"log"
	"github.com/tkanos/gonfig"
	"fmt"
	"github.com/kr/pretty"
)

type Configuration struct {

	// Filesystem

	BaseOutputDirectory			string	// Base directory where transient files are kept
	GeneratedModelDirectory		string	// Subdirectory where statistical models are kept
	PingResultDirectory			string	// Subdirectory where results of ping scans are kept
	NetworkGroupDirectory		string	// Subdirectory where results of grouping live hosts are kept
	NetworkBlacklistDirectory	string	// Subdirectory where network range blacklists are kept
	CleanPingResultDirectory	string	// Subdirectory where cleaned ping results are kept

	// Candidate address generation

	GenerateAddressCount		int		// How many addresses to generate in a given iteration

	// Network grouping and validation

	NetworkGroupingSize			int		// The bit-length of network size to use when checking for many-to-one
	NetworkPingCount			int		// The number of addresses to try pinging when testing for many-to-one

}

func LoadFromFile(filePath string) (Configuration, error) {
	log.Printf("Attempting to load config file from path '%s'.", filePath)
	config := Configuration{}
	err := gonfig.GetConf(filePath, &config)
	if err != nil {
		log.Printf("Error thrown when attempting to read config file at path '%s': %e", filePath, err)
		return Configuration{}, err
	} else {
		log.Printf("Successfully loaded config file from path '%s'.", filePath)
		return config, nil
	}
}

func (config *Configuration) Print() () {
	fmt.Print("\n-= Configuration Values =-\n\n")
	fmt.Printf("%# v", pretty.Formatter(config))
}
