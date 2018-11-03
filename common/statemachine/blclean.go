package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"log"
	"github.com/lavalamp-/ipv666/common/addressing"
	"github.com/rcrowley/go-metrics"
	"time"
	"github.com/lavalamp-/ipv666/common/fs"
)

var blRemovalDurationTimer = metrics.NewTimer()
var blRemovalCount = metrics.NewCounter()
var blLegitimateCount = metrics.NewCounter()

func init() {
	metrics.Register("blclean.removal.time", blRemovalDurationTimer)
	metrics.Register("blclean.removal.count", blRemovalCount)
	metrics.Register("blclean.legitimate.count", blLegitimateCount)
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
	cleanedAddrs := blacklist.CleanIPList(addrs, conf.LogLoopEmitFreq)
	elapsed := time.Since(start)
	blRemovalDurationTimer.Update(elapsed)
	blRemovalCount.Inc(int64(len(addrs) - len(cleanedAddrs)))
	blLegitimateCount.Inc(int64(len(cleanedAddrs)))
	log.Printf("Resulting cleaned list contains %d addresses (down from %d). Cleaned in %s.", len(cleanedAddrs), len(addrs), elapsed)
	outputPath := fs.GetTimedFilePath(conf.GetCleanPingDirPath())
	log.Printf("Writing resulting cleaned ping addresses to file at path '%s'.", outputPath)
	err = addressing.WriteIPsToBinaryFile(outputPath, cleanedAddrs)
	if err != nil {
		return err
	}
	log.Printf("Cleaned ping results successfully written to path '%s'.", outputPath)
	//TODO aggregate all found IP addresses
	data.UpdateCleanPingResults(cleanedAddrs, outputPath)
	return nil
}
