package statemachine

import (
	"github.com/lavalamp-/ipv666/common/blacklist"
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"github.com/lavalamp-/ipv666/common/fs"
	"github.com/rcrowley/go-metrics"
	"github.com/spf13/viper"
	"log"
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

	log.Print("Processing the aliased networks that were found into blacklist.")

	curBlacklist, err := data.GetBlacklist(config.GetNetworkBlacklistDirPath())
	if err != nil {
		return err
	}
	aliasedNets, err := data.GetAliasedNetworks()
	if err != nil {
		return err
	}

	log.Print("Loaded all relevant data into memory. Processing aliased results now.")

	start := time.Now()
	added, skipped := curBlacklist.AddNetworks(aliasedNets)
	elapsed := time.Since(start)
	aliasProcessTime.Update(elapsed)
	aliasProcessSkippedCount.Inc(int64(skipped))
	aliasProcessAddedCount.Inc(int64(added))

	log.Printf("Successfully processed %d aliased networks in %s. %d were added, %d were skipped.", len(aliasedNets), elapsed, added, skipped)

	log.Printf("Cleaning blacklist now. Blacklist is starting at capacity %d.", curBlacklist.GetCount())
	start = time.Now()
	numCleaned := curBlacklist.Clean(viper.GetInt("LogLoopEmitFreq"))
	aliasBlacklistCleanTime.Update(time.Since(start))
	aliasBlacklistCleanCount.Inc(int64(numCleaned))
	log.Printf("%d networks were cleaned from the blacklist (down to %d capacity).", numCleaned, curBlacklist.GetCount())

	outputPath := fs.GetTimedFilePath(config.GetNetworkBlacklistDirPath())
	log.Printf("Writing new blacklist to file at path '%s'.", outputPath)
	start = time.Now()
	err = blacklist.WriteNetworkBlacklistToFile(outputPath, curBlacklist)
	if err != nil {
		log.Printf("Error thrown when writing blacklist to file '%s': %e", outputPath, err)
		return err
	}
	aliasBlacklistWriteTime.Update(time.Since(start))

	data.UpdateBlacklist(curBlacklist, outputPath)

	log.Print("Successfully updated blacklist based on the results of the aliased network checking.")

	return nil

}
