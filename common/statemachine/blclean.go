package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"log"
	"github.com/lavalamp-/ipv666/common/addresses"
	"os"
	"fmt"
)

func cleanBlacklistedAddresses(conf *config.Configuration) (error) {

	// Find the blacklist file path
	blacklistPath, err := data.GetMostRecentFilePathFromDir(conf.GetNetworkBlacklistDirPath())
	if err != nil {
		return err
	}

	// Load the blacklist network addresses
	log.Printf("Loading blacklist network addresses")
	nets, err := addresses.GetAddressListFromHexStringsFile(blacklistPath)
	if err != nil {
		return err
	}

	// Find the ping results file path
	addrsPath, err := data.GetMostRecentFilePathFromDir(conf.GetPingResultDirPath())
	if err != nil {
		return err
	}

	// Load the ping results
	log.Printf("Loading ping scan result addresses")
	addrs, err := addresses.GetAddressListFromHexStringsFile(addrsPath)
	if err != nil {
		return err
	}

	// Remove addresses from blacklisted networks
	log.Printf("Removing addresses from blacklisted networks")
	var cleanAddrs []addresses.IPv6Address
	for _, addr := range(addrs.Addresses) {
		found := false
		for _, net := range(nets.Addresses) {
			match := true
			for x := 0; x < conf.NetworkGroupingSize; x++ {
				byteOff := (int)(x/8)
				bitOff := (uint)(x-(byteOff*8))
				byteMask := (byte)(1 << bitOff)
				if (addr.Content[byteOff] & byteMask) != (net.Content[byteOff] & byteMask) {
					match = false
					break
				}
			}
			if match == true {
				found = true
				break
			}
		}
		if found == false {
			cleanAddrs = append(cleanAddrs, addr)
		}
	}

	// Write the clean ping response addresses to disk
	cleanPath := getTimedFilePath(conf.GetCleanPingDirPath())
	log.Printf("Writing clean addresses to %s.", cleanPath)
	file, err := os.OpenFile(cleanPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	for _, addr := range(cleanAddrs) {
		file.WriteString(fmt.Sprintf("%s\n", addr.String()))
	}
	file.Close()
	return nil
}