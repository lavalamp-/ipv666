package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"log"
	"github.com/lavalamp-/ipv666/common/addressing"
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
	blacklist, err := data.GetBlacklist(conf.GetNetworkBlacklistDirPath())
	if err != nil {
		return err
	}
	log.Printf("Cleaning addresses using blacklist with %d entries.", len(blacklist.Networks))
	addrs, err := data.GetCandidatePingResults(conf.GetPingResultDirPath())
	if err != nil {
		return err
	}
	log.Printf("Total of %d addresses to clean.", len(addrs))
	start := time.Now()
	cleaned_addrs := blacklist.CleanIPList(addrs)
	elapsed := time.Since(start)
	blRemovalDurationTimer.Update(elapsed)
	blRemovalCountGauge.Update(int64(len(addrs) - len(cleaned_addrs)))
	blLegitimateCountGauge.Update(int64(len(cleaned_addrs)))
	log.Printf("Resulting cleaned list contains %d addresses (down from %d). Cleaned in %s.", len(cleaned_addrs), len(addrs), elapsed)
	outputPath := getTimedFilePath(conf.GetCleanPingDirPath())
	log.Printf("Writing resulting cleaned ping addresses to file at path '%s'.", outputPath)
	err = addressing.WriteIPsToBinaryFile(outputPath, cleaned_addrs)
	if err != nil {
		return err
	}
	log.Printf("Cleaned ping results successfully written to path '%s'.", outputPath)
	//TODO aggregate all found IP addresses
	data.UpdateCleanPingResults(cleaned_addrs, outputPath)
	return nil
}

//func cleanBlacklistedAddresses(conf *config.Configuration) (error) {
//
//	// Find the blacklist file path
//	blacklistPath, err := data.GetMostRecentFilePathFromDir(conf.GetNetworkBlacklistDirPath())
//	if err != nil {
//		return err
//	}
//
//	// Load the blacklist network addressing
//	log.Printf("Loading blacklist network addressing")
//	nets, err := addressing.GetAddressListFromHexStringsFile(blacklistPath)
//	if err != nil {
//		return err
//	}
//
//	// Find the ping results file path
//	addrsPath, err := data.GetMostRecentFilePathFromDir(conf.GetPingResultDirPath())
//	if err != nil {
//		return err
//	}
//
//	// Load the ping results
//	log.Printf("Loading ping scan result addressing")
//	addrs, err := addressing.GetAddressListFromHexStringsFile(addrsPath)
//	if err != nil {
//		return err
//	}
//
//	start := time.Now()
//	// Remove addressing from blacklisted networks
//	log.Printf("Removing addressing from blacklisted networks")
//	var cleanAddrs []addressing.IPv6Address
//	for _, addr := range(addrs.Addresses) {
//		found := false
//		for _, net := range(nets.Addresses) {
//			match := true
//			for x := 0; x < conf.NetworkGroupingSize; x++ {
//				byteOff := (int)(x/8)
//				bitOff := (uint)(x-(byteOff*8))
//				byteMask := (byte)(1 << bitOff)
//				if (addr.Content[byteOff] & byteMask) != (net.Content[byteOff] & byteMask) {
//					match = false
//					break
//				}
//			}
//			if match == true {
//				found = true
//				break
//			}
//		}
//		if found == false {
//			cleanAddrs = append(cleanAddrs, addr)
//		}
//	}
//	elapsed := time.Since(start)
//	blRemovalDurationTimer.Update(elapsed)
//	blRemovalCountGauge.Update(int64(len(addrs.Addresses) - len(cleanAddrs)))
//	blLegitimateCountGauge.Update(int64(len(cleanAddrs)))
//
//	// Write the clean ping response addressing to disk
//	cleanPath := getTimedFilePath(conf.GetCleanPingDirPath())
//	log.Printf("Writing clean addressing to %s.", cleanPath)
//	file, err := os.OpenFile(cleanPath, os.O_WRONLY|os.O_CREATE, 0600)
//	if err != nil {
//		return err
//	}
//	for _, addr := range(cleanAddrs) {
//		file.WriteString(fmt.Sprintf("%s\n", addr.String()))
//	}
//	file.Close()
//	return nil
//}