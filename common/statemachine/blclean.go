package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"log"
	"github.com/lavalamp-/ipv666/common/addresses"
	"os"
	"fmt"
	"github.com/rcrowley/go-metrics"
	"time"
)

var blRemovalDurationTimer = metrics.NewTimer()
var blRemovalCountGauge = metrics.NewGauge()
var blLegitimateCountGauge = metrics.NewGauge()

func init() {
	metrics.Register("bl_removal_duration", blRemovalDurationTimer)
	metrics.Register("bl_removal_count", blRemovalCountGauge)
	metrics.Register("bl_legitimate_count", blLegitimateCountGauge)
}

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

	start := time.Now()
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
	elapsed := time.Since(start)
	blRemovalDurationTimer.Update(elapsed)
	blRemovalCountGauge.Update(int64(len(addrs.Addresses) - len(cleanAddrs)))
	blLegitimateCountGauge.Update(int64(len(cleanAddrs)))

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