package config

import (
	"log"
	"github.com/tkanos/gonfig"
	"fmt"
	"github.com/kr/pretty"
	"path/filepath"
	"time"
)

type Configuration struct {

	// Filesystem

	BaseOutputDirectory			string	// Base directory where transient files are kept
	GeneratedModelDirectory		string	// Subdirectory where statistical models are kept
	CandidateAddressDirectory	string	// Subdirectory where generated candidate addressing are kept
	PingResultDirectory			string	// Subdirectory where results of ping scans are kept
	NetworkGroupDirectory		string	// Subdirectory where results of grouping live hosts are kept
	NetworkScanTargetsDirectory	string	// Subdirectory where the addresses to scan for blacklist checks are kept
	NetworkScanResultsDirectory	string	// Subdirectory where the results of scanning blacklist candidate networks are kept
	NetworkBlacklistDirectory	string	// Subdirectory where network range blacklists are kept
	CleanPingResultDirectory	string	// Subdirectory where cleaned ping results are kept
	BloomFilterDirectory		string	// Subdirectory where the Bloom filter is kept
	StateFileName				string	// The file name for the file that contains the current state

	// Candidate address generation

	GenerateAddressCount		int		// How many addressing to generate in a given iteration
	GenerateFirstNybble			uint8	// The first nybble of IPv6 addressing to generate
	GenerateUpdateFreq			int		// The interval upon which to emit to a log file during address generation
	GenWriteUpdateFreq			int		// The interval upon which to emit to a log file during writing address files
	ModelUpdateFreq				int		// The interval upon which to emit to a log file during model updates

	// Existing address bloom filter

	AddressFilterSize			uint	// The size of the Bloom filter to use for identifying already guessed addresses
	AddressFilterHashCount		uint	// The number of hashing functions to use for the address Bloom filter
	BloomEmptyMultiple			float64	// The multiple of the address generation size upon which the Bloom filter should be emptied and remade

	// Network grouping and validation

	NetworkGroupingSize			uint8	// The bit-length of network size to use when checking for many-to-one
	NetworkPingCount			int		// The number of addressing to try pinging when testing for many-to-one
	NetworkBlacklistPercent		float64	// The percentage of ping results that, if returned positive, indicate a blacklisted network

	// Blacklist candidate generation

	BlacklistFlushInterval		int		// The frequency with which to write newly-generate blacklist candidate addresses to disk

	// Logging

	LogToFile					bool	// Whether or not to write log results to a file instead of stdout
	LogFilePath					string	// The local file path to where log files should be written
	LogFileMBSize				int		// The max size of each log file in MB
	LogFileMaxBackups			int		// The maximum number of backups to have in rotating log files
	LogFileMaxAge				int		// The maximum number of days to store log files
	CompressLogFiles			bool	// Whether or not to compress log files
	LogLoopEmitFreq				int		// The general frequency with which logs should be emitted in long loops

	// Scanning

	ZmapExecPath				string  // Local file path to the Zmap executable
	ZmapBandwidth				string  // Bandwidth cap for Zmap
	ZmapSourceAddress   		string  // Source IPv6 address for Zmap

	// Exportation

	ExportEnabled				bool	// Whether or not to export data to S3
	ExitOnFailedSync			bool	// Whether or not to exit the program if an S3 sync fails

	// AWS

	AWSBucketRegion				string	// The region where the AWS S3 bucket resides
	AWSBucketName				string	// The name of the bucket to push to
	AWSAccessKey				string	// The AWS access key to use
	AWSSecretKey				string	// The AWS secret key to use

	// Clean Up

	CleanUpEnabled				bool	// Whether or not to delete non-recent files after a run

	// Metrics

	MetricsStateLoopPrefix		string	// The prefix for the state loop metric
	ExitOnFailedMetrics			bool	// Whether or not to exit the program when a metrics operation fails
	MetricsToStdout				bool	// Whether or not to print metrics to Stdout
	MetricsStdoutFreq			int64	// The frequency in seconds of how often to print metrics to Stdout
	GraphiteExportEnabled		bool	// Whether or not to export data to Graphite
	GraphiteHost				string	// The host address for Graphite
	GraphitePort				int		// The Graphite port
	GraphiteEmitFreq			int64	// How often to emit metrics to Graphite in seconds

	// Output

	OutputFileName				string	// The file name for the file to write addresses to
	OutputFileType				string	// The output file type

	// Input

	InputEntropyThreshold		float64	// The threshold upon which addresses having more entropy will be removed
	InputEntropyBitLength		int		// The number of bits within IP addresses to calculate entropy based on
	InputMinAddresses			int		// The recommended minimum number of addresses to require for a given statistical model
	InputEmitFreq				int		// The interval upon which to emit log updates when processing input files

	// Runtime

	ForceAcceptPrompts			bool	// Whether or not to bypass prompts by force accepting them

	// Alias Detection

	AliasLeftIndexStart			uint8	// The left-most index for CIDR mask lengths where aliased network detection should start
	AliasDuplicateScanCount		uint8	// The number of times a single address should be scanned when checking for aliased networks

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

func (config *Configuration) Print() {
	fmt.Print("\n-= Configuration Values =-\n\n")
	fmt.Printf("%# v", pretty.Formatter(config))
}

func (config *Configuration) GetOutputFilePath() (string) {
	return fmt.Sprintf("%s.%s", config.OutputFileName, config.OutputFileType)
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

func (config *Configuration) GetNetworkScanTargetsDirPath() (string) {
	return filepath.Join(config.BaseOutputDirectory, config.NetworkScanTargetsDirectory)
}

func (config *Configuration) GetNetworkScanResultsDirPath() (string) {
	return filepath.Join(config.BaseOutputDirectory, config.NetworkScanResultsDirectory)
}

func (config *Configuration) GetNetworkBlacklistDirPath() (string) {
	return filepath.Join(config.BaseOutputDirectory, config.NetworkBlacklistDirectory)
}

func (config *Configuration) GetCleanPingDirPath() (string) {
	return filepath.Join(config.BaseOutputDirectory, config.CleanPingResultDirectory)
}

func (config *Configuration) GetBloomDirPath() (string) {
	return filepath.Join(config.BaseOutputDirectory, config.BloomFilterDirectory)
}

func (config *Configuration) GetAllDirectories() ([]string) {
	return []string{
		config.BaseOutputDirectory,
		config.GetGeneratedModelDirPath(),
		config.GetCandidateAddressDirPath(),
		config.GetPingResultDirPath(),
		config.GetNetworkGroupDirPath(),
		config.GetNetworkScanTargetsDirPath(),
		config.GetNetworkScanResultsDirPath(),
		config.GetNetworkBlacklistDirPath(),
		config.GetCleanPingDirPath(),
		config.GetBloomDirPath(),
	}
}

func (config *Configuration) GetAllExportDirectories() ([]string) {
	return []string{
		config.GetGeneratedModelDirPath(),
		config.GetCandidateAddressDirPath(),
		config.GetPingResultDirPath(),
		config.GetNetworkGroupDirPath(),
		config.GetNetworkScanTargetsDirPath(),
		config.GetNetworkScanResultsDirPath(),
		config.GetNetworkBlacklistDirPath(),
		config.GetCleanPingDirPath(),
		config.GetBloomDirPath(),
	}
}

func (config *Configuration) GetGraphiteEmitDuration() (time.Duration) {
	return time.Duration(config.GraphiteEmitFreq) * time.Second
}