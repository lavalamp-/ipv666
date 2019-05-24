package config

import (
	"fmt"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"net"
	"path"
	"path/filepath"
	"time"
)

func InitConfig() {
	viper.SetEnvPrefix("IPV666")

	// Filesystem

	viper.BindEnv("BaseOutputDirectory")				// Base directory where transient files are kept
	viper.BindEnv("GeneratedModelDirectory")			// Subdirectory where statistical models are kept
	viper.BindEnv("CandidateAddressDirectory")		// Subdirectory where generated candidate addressing are kept
	viper.BindEnv("PingResultDirectory")				// Subdirectory where results of ping scans are kept
	viper.BindEnv("NetworkGroupDirectory")			// Subdirectory where results of grouping live hosts are kept
	viper.BindEnv("NetworkScanTargetsDirectory")		// Subdirectory where the addresses to scan for blacklist checks are kept
	viper.BindEnv("NetworkScanResultsDirectory")		// Subdirectory where the results of scanning blacklist candidate networks are kept
	viper.BindEnv("NetworkBlacklistDirectory")		// Subdirectory where network range blacklists are kept
	viper.BindEnv("CleanPingResultDirectory")		// Subdirectory where cleaned ping results are kept
	viper.BindEnv("AliasedNetworkDirectory")			// Subdirectory where aliased network results are kept
	viper.BindEnv("BloomFilterDirectory")			// Subdirectory where the Bloom filter is kept
	viper.BindEnv("StateFileName")					// The file name for the file that contains the current state
	viper.BindEnv("TargetNetworkFileName")			// The file name for the file that contains the last network that was targeted
	viper.BindEnv("CloudSyncOptInPath")			// Cloud sync opt-in status file path
	viper.BindEnv("CloudSyncOptIn")			// Cloud sync opt-in status

	home, err := homedir.Dir()
	if err != nil {
		logging.ErrorF(err)
	}

	viper.SetDefault("BaseOutputDirectory", path.Join(home, ".ipv666"))
	viper.SetDefault("GeneratedModelDirectory", "models")
	viper.SetDefault("CandidateAddressDirectory", "candidates")
	viper.SetDefault("PingResultDirectory", "pingresult")
	viper.SetDefault("NetworkGroupDirectory", "networkgroups")
	viper.SetDefault("NetworkScanTargetsDirectory", "networkscantargets")
	viper.SetDefault("NetworkScanResultsDirectory", "networkscanresults")
	viper.SetDefault("NetworkBlacklistDirectory", "networkblacklist")
	viper.SetDefault("CleanPingResultDirectory", "cleanpings")
	viper.SetDefault("AliasedNetworkDirectory", "aliasednets")
	viper.SetDefault("BloomFilterDirectory", "bloom")
	viper.SetDefault("StateFileName", "state.bin")
	viper.SetDefault("TargetNetworkFileName", "network.bin")
	viper.SetDefault("CloudSyncOptInPath", ".cloudsyncoptin")
	viper.SetDefault("CloudSyncOptIn", false)

	// Candidate address generation

	viper.BindEnv("GenerateAddressCount")			// How many addressing to generate in a given iteration

	viper.SetDefault("GenerateAddressCount", 1000000)

	// Modeling

	viper.BindEnv("ModelGenerationJitter")			// The default jitter (ie: % likelihood of a random wildcard) to use when generating new addresses
	viper.BindEnv("ModelCheckCount")					// The number of upgrades to wait between checking for cluster model improvement
	viper.BindEnv("ModelMinNybblePercent")			// The minimum percent probability of a nybble occurring in a cluster model
	viper.BindEnv("ModelDistributionSize")			// The size of startdust distributions used for random nybble generation

	viper.SetDefault("ModelGenerationJitter", 0.1) //TODO figure out why configuration stuff isn't working with cobra's bindpflags
	viper.SetDefault("ModelCheckCount", 10000)
	viper.SetDefault("ModelMinNybblePercent", 0.01)
	viper.SetDefault("ModelDistributionSize", 1000)

	// Existing address bloom filter

	viper.BindEnv("AddressFilterSize")				// The size of the Bloom filter to use for identifying already guessed addresses
	viper.BindEnv("AddressFilterHashCount")			// The number of hashing functions to use for the address Bloom filter
	viper.BindEnv("BloomEmptyMultiple")				// The multiple of the address generation size upon which the Bloom filter should be emptied and remade

	viper.SetDefault("AddressFilterSize", 250000000)
	viper.SetDefault("AddressFilterHashCount", 3)
	viper.SetDefault("BloomEmptyMultiple", 2.0)

	// Network grouping and validation

	viper.BindEnv("NetworkGroupingSize")				// The bit-length of network size to use when checking for many-to-one
	viper.BindEnv("NetworkPingCount")				// The number of addressing to try pinging when testing for many-to-one
	viper.BindEnv("NetworkBlacklistPercent")			// The percentage of ping results that, if returned positive, indicate a blacklisted network

	viper.SetDefault("NetworkGroupingSize", 96)
	viper.SetDefault("NetworkPingCount", 6)
	viper.SetDefault("NetworkBlacklistPercent", 0.5)

	// Blacklist candidate generation

	viper.BindEnv("BlacklistFlushInterval")			// The frequency with which to write newly-generate blacklist candidate addresses to disk

	viper.SetDefault("BlacklistFlushInterval", 500000)

	// Fan-out ping-scanning
	viper.BindEnv("FanOutNetworkBlockSize")        // Number of contiguous neighboring /64 networks to attempt
	viper.BindEnv("FanOutHostBlockSize")           // Number of contiguous hosts to attempt, monotonically increasing from each /64
	viper.BindEnv("FanOutMaxNetworks")             // Maximum networks to attempt during fan-out scanning
	viper.BindEnv("FanOutMaxHosts")                // Maximum hosts to attempt during fan-out scanning
	viper.SetDefault("FanOutNetworkBlockSize", 1000)
	viper.SetDefault("FanOutHostBlockSize", 500)
	viper.SetDefault("FanOutMaxNetworks", 2000000)
	viper.SetDefault("FanOutMaxHosts", 2000000)

	// Logging

	viper.BindEnv("LogLevel")						// The level to log at (debug, info, success, warn, error)
	viper.BindEnv("LogToFile")						// Whether or not to write log results to a file instead of stdout
	viper.BindEnv("LogFilePath")						// The local file path to where log files should be written
	viper.BindEnv("LogFileMBSize")					// The max size of each log file in MB
	viper.BindEnv("LogFileMaxBackups")				// The maximum number of backups to have in rotating log files
	viper.BindEnv("LogFileMaxAge")					// The maximum number of days to store log files
	viper.BindEnv("CompressLogFiles")				// Whether or not to compress log files
	viper.BindEnv("LogLoopEmitFreq")					// The general frequency with which logs should be emitted in long loops

	viper.SetDefault("LogLevel", "info")
	viper.SetDefault("LogToFile", false)
	viper.SetDefault("LogFilePath", "ipv666.log")
	viper.SetDefault("LogFileMBSize", 10)
	viper.SetDefault("LogFileMaxBackups", 10)
	viper.SetDefault("LogFileMaxAge", 120)
	viper.SetDefault("CompressLogFiles", false)
	viper.SetDefault("LogLoopEmitFreq", 250000)

	// Scanning

	viper.BindEnv("PingScanBandwidth")				// The maximum bandwidth to use for ping scanning
	viper.BindEnv("ScanTargetNetwork")				// The default network to scan

	viper.SetDefault("PingScanBandwidth", "20M")
	viper.SetDefault("ScanTargetNetwork", "2000::/4")

	// Clean Up

	viper.BindEnv("CleanUpEnabled")					// Whether or not to delete non-recent files after a run

	viper.SetDefault("CleanUpEnabled", true)

	// Metrics

	viper.BindEnv("ExitOnFailedMetrics")				// Whether or not to exit the program when a metrics operation fails
	viper.BindEnv("MetricsToStdout")					// Whether or not to print metrics to Stdout
	viper.BindEnv("MetricsStdoutFreq")				// The frequency in seconds of how often to print metrics to Stdout
	viper.BindEnv("GraphiteExportEnabled")			// Whether or not to export data to Graphite
	viper.BindEnv("GraphiteHost")					// The host address for Graphite
	viper.BindEnv("GraphitePort")					// The Graphite port
	viper.BindEnv("GraphiteEmitFreq")				// How often to emit metrics to Graphite in seconds

	viper.SetDefault("ExitOnFailedMetrics", false)
	viper.SetDefault("MetricsToStdout", false)
	viper.SetDefault("MetricsStdoutFreq", 300)
	viper.SetDefault("GraphiteExportEnabled", false)
	viper.SetDefault("GraphiteHost", "127.0.0.1")
	viper.SetDefault("GraphitePort", 2003)
	viper.SetDefault("GraphiteEmitFreq", 60)

	// Output

	viper.BindEnv("OutputFileName")					// The file name for the file to write addresses to
	viper.BindEnv("OutputFileType")					// The output file type

	viper.SetDefault("OutputFileName", "discovered_addrs")  //TODO remove default output file name and type
	viper.SetDefault("OutputFileType", "txt")

	// Input

	viper.BindEnv("InputMinTargetCount")				// The minimum bit count for network sizes to scan

	viper.SetDefault("InputMinTargetCount", 30)

	// Runtime

	viper.BindEnv("ForceAcceptPrompts")				// Whether or not to bypass prompts by force accepting them

	viper.SetDefault("ForceAcceptPrompts", false)

	// Alias Detection

	viper.BindEnv("AliasLeftIndexStart")				// The left-most index for CIDR mask lengths where aliased network detection should start
	viper.BindEnv("AliasDuplicateScanCount")			// The number of times a single address should be scanned when checking for aliased networks

	viper.SetDefault("AliasLeftIndexStart", 0)
	viper.SetDefault("AliasDuplicateScanCount", 3)

	viper.AutomaticEnv()
}

func GetCloudSyncOptInPath() string {
	return filepath.Join(viper.GetString("BaseOutputDirectory"), viper.GetString("CloudSyncOptInPath"))
}

func GetCloudSyncOptIn() bool {
	return viper.GetBool("CloudSyncOptIn")
}

func SetCloudSyncOptIn(state bool) {
	viper.Set("CloudSyncOptIn", state)
}

func GetOutputFilePath() string {
	return fmt.Sprintf("%s.%s", viper.GetString("OutputFileName"), viper.GetString("OutputFileType"))
}

func GetStateFilePath() string {
	return filepath.Join(viper.GetString("BaseOutputDirectory"), viper.GetString("StateFileName"))
}

func GetTargetNetworkFilePath() string {
	return filepath.Join(viper.GetString("BaseOutputDirectory"), viper.GetString("TargetNetworkFileName"))
}

func GetGeneratedModelDirPath() string {
	return filepath.Join(viper.GetString("BaseOutputDirectory"), viper.GetString("GeneratedModelDirectory"))
}

func GetCandidateAddressDirPath() string {
	return filepath.Join(viper.GetString("BaseOutputDirectory"), viper.GetString("CandidateAddressDirectory"))
}

func GetPingResultDirPath() string {
	return filepath.Join(viper.GetString("BaseOutputDirectory"), viper.GetString("PingResultDirectory"))
}

func GetNetworkGroupDirPath() string {
	return filepath.Join(viper.GetString("BaseOutputDirectory"), viper.GetString("NetworkGroupDirectory"))
}

func GetNetworkScanTargetsDirPath() string {
	return filepath.Join(viper.GetString("BaseOutputDirectory"), viper.GetString("NetworkScanTargetsDirectory"))
}

func GetNetworkScanResultsDirPath() string {
	return filepath.Join(viper.GetString("BaseOutputDirectory"), viper.GetString("NetworkScanResultsDirectory"))
}

func GetNetworkBlacklistDirPath() string {
	return filepath.Join(viper.GetString("BaseOutputDirectory"), viper.GetString("NetworkBlacklistDirectory"))
}

func GetCleanPingDirPath() string {
	return filepath.Join(viper.GetString("BaseOutputDirectory"), viper.GetString("CleanPingResultDirectory"))
}

func GetAliasedNetworkDirPath() string {
	return filepath.Join(viper.GetString("BaseOutputDirectory"), viper.GetString("AliasedNetworkDirectory"))
}

func GetBloomDirPath() string {
	return filepath.Join(viper.GetString("BaseOutputDirectory"), viper.GetString("BloomFilterDirectory"))
}

func GetAllDirectories() []string {
	return []string{
		viper.GetString("BaseOutputDirectory"),
		GetGeneratedModelDirPath(),
		GetCandidateAddressDirPath(),
		GetPingResultDirPath(),
		GetNetworkGroupDirPath(),
		GetNetworkScanTargetsDirPath(),
		GetNetworkScanResultsDirPath(),
		GetNetworkBlacklistDirPath(),
		GetCleanPingDirPath(),
		GetAliasedNetworkDirPath(),
		GetBloomDirPath(),
	}
}

func GetAllExportDirectories() []string {
	return []string{
		GetGeneratedModelDirPath(),
		GetCandidateAddressDirPath(),
		GetPingResultDirPath(),
		GetNetworkGroupDirPath(),
		GetNetworkScanTargetsDirPath(),
		GetNetworkScanResultsDirPath(),
		GetNetworkBlacklistDirPath(),
		GetCleanPingDirPath(),
		GetAliasedNetworkDirPath(),
		GetBloomDirPath(),
	}
}

func GetGraphiteEmitDuration() time.Duration {
	return time.Duration(viper.GetInt64("GraphiteEmitFreq")) * time.Second
}

func GetTargetNetwork() (*net.IPNet, error) {
	_, network, err := net.ParseCIDR(viper.GetString("ScanTargetNetwork"))
	return network, err
}
