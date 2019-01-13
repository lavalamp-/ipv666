package validation

import (
	"errors"
	"fmt"
	"github.com/lavalamp-/ipv666/common/addressing"
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"github.com/lavalamp-/ipv666/common/fs"
	"github.com/spf13/viper"
	"net"
)

func ValidateIPv6NetworkString(toParse string) (*net.IPNet, error) {
	ip, targetNetwork, err := net.ParseCIDR(toParse)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error thrown when parsing target network string of '%s': %e", toParse, err))
	} else if targetNetwork == nil {
		return nil, errors.New(fmt.Sprintf("Target network of '%s' could not be decoded to a valid IPv6 network range.", toParse))
	} else if addressing.IsAddressIPv4(&ip) {
		return nil, errors.New(fmt.Sprintf("Network range '%s' decoded to an IPv4 network range (%s).", toParse, ip))
	}
	return targetNetwork, nil
}

func ValidateIPv6NetworkStringForScanning(toParse string) (*net.IPNet, error) {
	//TODO check to see if value is in any of the weird predefined network ranges
	network, err := ValidateIPv6NetworkString(toParse)
	if err != nil {
		return nil, err
	}
	curBlacklist, err := data.GetBlacklist()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error thrown when reading blacklist from directory '%s': %e", config.GetNetworkBlacklistDirPath(), err))
	}
	if curBlacklist.IsNetworkBlacklisted(network) {
		blacklistingNetwork := curBlacklist.GetBlacklistingNetworkFromNetwork(network)
		return nil, errors.New(fmt.Sprintf("The network you picked (%s) is already blacklisted (by network %s), indicating that the network is aliased. Scanning this network will not result in any actionable information.", network, blacklistingNetwork))
	}
	ones, _ := network.Mask.Size()
	bitsLeft := 128 - ones
	if bitsLeft < viper.GetInt("InputMinTargetCount") {
		return nil, errors.New(fmt.Sprintf("You specified a network range that had 2^30 or less addresses in it (specifically 2^%d). This tool is not meant for such small ranges. We recommend using the IPv6 Zmap directly for this whole range.", bitsLeft))
	}
	return network, nil
}

func ValidateOutputFileType(toCheck string) error {
	if toCheck == "txt" || toCheck == "bin" || toCheck == "hex" {
		return nil
	} else {
		return fmt.Errorf("%s is not a valid output file type (expected 'txt', 'bin', or 'hex')", toCheck)
	}
}

func ValidateLogLevel(toCheck string) error {
	if toCheck == "debug" || toCheck == "info" || toCheck == "success" || toCheck == "warning" || toCheck == "error" {
		return nil
	} else {
		return fmt.Errorf("'%s' is not a valid log level (expected one of 'debug', 'info', 'success', 'warning', or 'error')", toCheck)
	}
}

func ValidateFileNotExist(filePath string) error {
	if fs.CheckIfFileExists(filePath) {
		return fmt.Errorf("a file already exists at path '%s'", filePath)
	} else {
		return nil
	}
}

func ValidateFileExists(filePath string) error {
	if !fs.CheckIfFileExists(filePath) {
		return fmt.Errorf("no file found at path '%s'", filePath)
	} else {
		return nil
	}
}
