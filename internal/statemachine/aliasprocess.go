package statemachine

import (
	"github.com/lavalamp-/ipv666/internal/blacklist"
	"github.com/lavalamp-/ipv666/internal/config"
	"github.com/lavalamp-/ipv666/internal/data"
	"github.com/lavalamp-/ipv666/internal/fs"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/rcrowley/go-metrics"
	"github.com/spf13/viper"
	"time"
)

var aliasProcessAddedCount = metrics.NewCounter()
var aliasProcessSkippedCount = metrics.NewCounter()
var aliasProcessTime = metrics.NewTimer()
var aliasBlacklistWriteTime = metrics.NewTimer()
var aliasBlacklistCleanTime = metrics.NewTimer()
var aliasBlacklistCleanCount = metrics.NewCounter()

func init() {
	metrics.Register("aliasprocess.process.added.count", aliasProcessAddedCount)
	metrics.Register("aliasprocess.process.skipped.count", aliasProcessSkippedCount)
	metrics.Register("aliasprocess.process.time", aliasProcessTime)
	metrics.Register("aliasprocess.blacklist.write.time", aliasBlacklistWriteTime)
	metrics.Register("aliasprocess.blacklist.clean.time", aliasBlacklistCleanTime)
	metrics.Register("aliasprocess.blacklist.clean.count", aliasBlacklistCleanCount)
}

func processAliasedNetworks() error {

	logging.Infof("Processing the aliased networks that were found into blacklist.")

	curBlacklist, err := data.GetBlacklist()
	if err != nil {
		return err
	}
	aliasedNets, err := data.GetAliasedNetworks()
	if err != nil {
		return err
	}

	logging.Debugf("Loaded all relevant data into memory. Processing aliased results now.")

	start := time.Now()
	added, skipped := curBlacklist.AddNetworks(aliasedNets)
	elapsed := time.Since(start)
	aliasProcessTime.Update(elapsed)
	aliasProcessSkippedCount.Inc(int64(skipped))
	aliasProcessAddedCount.Inc(int64(added))

	logging.Debugf("Successfully processed %d aliased networks in %s. %d were added, %d were skipped.", len(aliasedNets), elapsed, added, skipped)

	logging.Debugf("Cleaning blacklist now. Blacklist is starting at capacity %d.", curBlacklist.GetCount())
	start = time.Now()
	numCleaned := curBlacklist.Clean(viper.GetInt("LogLoopEmitFreq"))
	aliasBlacklistCleanTime.Update(time.Since(start))
	aliasBlacklistCleanCount.Inc(int64(numCleaned))
	logging.Debugf("%d networks were cleaned from the blacklist (down to %d capacity).", numCleaned, curBlacklist.GetCount())

	outputPath := fs.GetTimedFilePath(config.GetNetworkBlacklistDirPath())
	logging.Debugf("Writing new blacklist to file at path '%s'.", outputPath)
	start = time.Now()
	err = blacklist.WriteNetworkBlacklistToFile(outputPath, curBlacklist)
	if err != nil {
		logging.Warnf("Error thrown when writing blacklist to file '%s': %e", outputPath, err)
		return err
	}
	aliasBlacklistWriteTime.Update(time.Since(start))

	data.UpdateBlacklist(curBlacklist, outputPath)

	logging.Infof("Successfully updated blacklist based on the results of the aliased network checking.")

	return nil

}
