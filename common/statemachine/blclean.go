package statemachine

import (
	"github.com/lavalamp-/ipv666/common/addressing"
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"github.com/lavalamp-/ipv666/common/fs"
	"github.com/lavalamp-/ipv666/common/logging"
	"github.com/rcrowley/go-metrics"
	"github.com/spf13/viper"
	"time"
)

var blRemovalDurationTimer = metrics.NewTimer()
var blRemovalCount = metrics.NewCounter()
var blLegitimateCount = metrics.NewCounter()

func init() {
	metrics.Register("blclean.removal.time", blRemovalDurationTimer)
	metrics.Register("blclean.removal.count", blRemovalCount)
	metrics.Register("blclean.legitimate.count", blLegitimateCount)
}

func cleanBlacklistedAddresses() error {
	blacklist, err := data.GetBlacklist(config.GetNetworkBlacklistDirPath())
	if err != nil {
		return err
	}
	logging.Infof("Cleaning addresses using blacklist with %d entries.", blacklist.GetCount())
	addrs, err := data.GetCandidatePingResults(config.GetPingResultDirPath())
	if err != nil {
		return err
	}
	logging.Debugf("Total of %d addresses to clean.", len(addrs))
	start := time.Now()
	cleanedAddrs := blacklist.CleanIPList(addrs, viper.GetInt("LogLoopEmitFreq"))
	elapsed := time.Since(start)
	blRemovalDurationTimer.Update(elapsed)
	blRemovalCount.Inc(int64(len(addrs) - len(cleanedAddrs)))
	blLegitimateCount.Inc(int64(len(cleanedAddrs)))
	logging.Debugf("Resulting cleaned list contains %d addresses (down from %d). Cleaned in %s.", len(cleanedAddrs), len(addrs), elapsed)
	outputPath := fs.GetTimedFilePath(config.GetCleanPingDirPath())
	logging.Debugf("Writing resulting cleaned ping addresses to file at path '%s'.", outputPath)
	err = addressing.WriteIPsToBinaryFile(outputPath, cleanedAddrs)
	if err != nil {
		return err
	}
	logging.Debugf("Cleaned ping results successfully written to path '%s'.", outputPath)
	data.UpdateCleanPingResults(cleanedAddrs, outputPath)
	return nil
}
