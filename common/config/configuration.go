package config

import (
	"log"
	"github.com/tkanos/gonfig"
	"fmt"
	"github.com/kr/pretty"
	"path/filepath"
)

type Configuration struct {

	// Filesystem

	BaseOutputDirectory			string	// Base directory where transient files are kept
	GeneratedModelDirectory		string	// Subdirectory where statistical models are kept
	CandidateAddressDirectory	string	// Subdirectory where generated candidate addresses are kept
	PingResultDirectory			string	// Subdirectory where results of ping scans are kept
	NetworkGroupDirectory		string	// Subdirectory where results of grouping live hosts are kept
	NetworkBlacklistDirectory	string	// Subdirectory where network range blacklists are kept
	CleanPingResultDirectory	string	// Subdirectory where cleaned ping results are kept
	StateFileName				string	// The file name for the file that contains the current state

	// Candidate address generation

	GenerateAddressCount		int		// How many addresses to generate in a given iteration
	GenerateFirstNybble			uint8	// The first nybble of IPv6 addresses to generate
	GenerateUpdateFreq			int		// The interval upon which to emit to a log file during address generation
	GenWriteUpdateFreq			int		// The interval upon which to emit to a log file during writing address files

	// Network grouping and validation

	NetworkGroupingSize			int		// The bit-length of network size to use when checking for many-to-one
	NetworkPingCount			int		// The number of addresses to try pinging when testing for many-to-one

	// Logging

	LogFilePath					string	// The local file path to where log files should be written
	LogFileMBSize				int		// The max size of each log file in MB
	LogFileMaxBackups			int		// The maximum number of backups to have in rotating log files
	LogFileMaxAge				int		// The maximum number of days to store log files
	CompressLogFiles			bool	// Whether or not to compress log files

	// Scanning

	ZmapExecPath				string  // Local file path to the Zmap executable
	ZmapBandwidth				string  // Bandwidth cap for Zmap

	// Exportation

	ExportEnabled				bool	// Whether or not to export data to S3

	// AWS

	AWSBucketRegion				string	// The region where the AWS S3 bucket resides
	AWSBucketName				string	// The name of the bucket to push to
	AWSAccessKey				string	// The AWS access key to use
	AWSSecretKey				string	// The AWS secret key to use

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

func (config *Configuration) GetStateFilePath() (string) {
	return filepath.Join(config.BaseOutputDirectory, config.StateFileName)
}

func (config *Configuration) GetGeneratedModelDirPath() (string) {
	return filepath.Join(config.BaseOutputDirectory, config.GeneratedModelDirectory)
}

func (config *Configuration) GetCandidateAddressDirPath() (string) {
	return filepath.Join(config.BaseOutputDirectory, config.CandidateAddressDirectory)
}

func (config *Configuration) GetPingResultDirPath() (string) {
	return filepath.Join(config.BaseOutputDirectory, config.PingResultDirectory)
}

func (config *Configuration) GetNetworkGroupDirPath() (string) {
	return filepath.Join(config.BaseOutputDirectory, config.NetworkGroupDirectory)
}

func (config *Configuration) GetNetworkBlacklistDirPath() (string) {
	return filepath.Join(config.BaseOutputDirectory, config.NetworkBlacklistDirectory)
}

func (config *Configuration) GetCleanPingDirPath() (string) {
	return filepath.Join(config.BaseOutputDirectory, config.CleanPingResultDirectory)
}

func (config *Configuration) GetAllDirectories() ([]string) {
	return []string{
		config.BaseOutputDirectory,
		config.GetGeneratedModelDirPath(),
		config.GetCandidateAddressDirPath(),
		config.GetPingResultDirPath(),
		config.GetNetworkGroupDirPath(),
		config.GetNetworkBlacklistDirPath(),
		config.GetCleanPingDirPath(),
	}
}

func (config *Configuration) GetAllExportDirectories() ([]string) {
	return []string{
		config.GetGeneratedModelDirPath(),
		config.GetCandidateAddressDirPath(),
		config.GetPingResultDirPath(),
		config.GetNetworkGroupDirPath(),
		config.GetNetworkBlacklistDirPath(),
		config.GetCleanPingDirPath(),
	}
}