package validation

import (
	"net"
	"github.com/lavalamp-/ipv666/common/addressing"
	"github.com/lavalamp-/ipv666/common/data"
	"github.com/lavalamp-/ipv666/common/config"
	"errors"
	"fmt"
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

func ValidateIPv6NetworkStringForScanning(toParse string, conf *config.Configuration) (*net.IPNet, error) {
	//TODO check to see if value is in any of the weird predefined network ranges
	network, err := ValidateIPv6NetworkString(toParse)
	if err != nil {
		return nil, err
	}
	curBlacklist, err := data.GetBlacklist(conf.GetNetworkBlacklistDirPath())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error thrown when reading blacklist from directory '%s': %e", conf.GetNetworkBlacklistDirPath(), err))
	}
	if curBlacklist.IsNetworkBlacklisted(network) {
		blacklistingNetwork := curBlacklist.GetBlacklistingNetworkFromNetwork(network)
		return nil, errors.New(fmt.Sprintf("The network you picked (%s) is already blacklisted (by network %s), indicating that the network is aliased. Scanning this network will not result in any actionable information.", network, blacklistingNetwork))
	}
	ones, _ := network.Mask.Size()
	bitsLeft := 128 - ones
	if bitsLeft < conf.InputMinTargetCount {
		return nil, errors.New(fmt.Sprintf("You specified a network range that had 2^30 or less addresses in it (specifically 2^%d). This tool is not meant for such small ranges. We recommend using the IPv6 Zmap directly for this whole range.", bitsLeft))
	}
	return network, nil
}
