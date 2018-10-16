package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"log"
	"time"
	"net"
	"github.com/rcrowley/go-metrics"
	"github.com/lavalamp-/ipv666/common/addressing"

	"github.com/lavalamp-/ipv666/common/fs"
	"github.com/lavalamp-/ipv666/common/filtering"
)

var generateDurationTimer = metrics.NewTimer()
var generateBlacklistCount = metrics.NewCounter()
var generateBloomCount = metrics.NewCounter()
var generateWriteTimer = metrics.NewTimer()
var bloomWriteTimer = metrics.NewTimer()

func init() {
	metrics.Register("gen_duration_timer", generateDurationTimer)
	metrics.Register("gen_blacklist_count", generateBlacklistCount)
	metrics.Register("gen_bloom_count", generateBloomCount)
	metrics.Register("gen_write_timer", generateWriteTimer)
	metrics.Register("bloom_write_timer", bloomWriteTimer)
}

func generateCandidateAddresses(conf *config.Configuration) (error) {

	// Load the statistical model, blacklist, and bloom filter

	model, err := data.GetProbabilisticAddressModel(conf.GetGeneratedModelDirPath())
	if err != nil {
		return err
	}
	blacklist, err := data.GetBlacklist(conf.GetNetworkBlacklistDirPath())
	if err != nil {
		return err
	}
	bloom, err := data.GetBloomFilter(conf)
	if err != nil {
		return err
	}

	// Generate all of the addresses and filter out based on Bloom filter and blacklist

	log.Printf(
		"Generating a total of %d addresses based on the content of model '%s' (%d digest count). Starting nybble is %d.",
		conf.GenerateAddressCount,
		model.Name,
		model.DigestCount,
		conf.GenerateFirstNybble,
	)
	var addresses []*net.IP
	var blacklistCount, bloomCount, madeCount = 0, 0, 0
	start := time.Now()
	for len(addresses) < conf.GenerateAddressCount {
		newIP := model.GenerateSingleIP(conf.GenerateFirstNybble)
		ipBytes := ([]byte)(*newIP)
		if blacklist.IsIPBlacklisted(newIP) {
			blacklistCount++
		} else if bloom.Test(ipBytes) {
			bloomCount++
		} else {
			madeCount++
			addresses = append(addresses, newIP)
			bloom.Add(ipBytes)
		}
		if (madeCount + blacklistCount + bloomCount) % conf.GenerateUpdateFreq == 0 {
			log.Printf("Generated %d total addresses, %d have been valid, %d have been blacklisted, %d exist in Bloom filter.", madeCount + blacklistCount, madeCount, blacklistCount, bloomCount)
		}
	}
	elapsed := time.Since(start)
	generateDurationTimer.Update(elapsed)
	generateBlacklistCount.Inc(int64(blacklistCount))
	generateBloomCount.Inc(int64(bloomCount))
	log.Printf("Took a total of %s to generate %d candidate addresses (%d blacklisted filtered out, %d existed in Bloom filter).", elapsed, conf.GenerateAddressCount, blacklistCount, bloomCount)

	// Write addresses and Bloom filter to disk and update data manager to point to in-memory references

	outputPath := fs.GetTimedFilePath(conf.GetCandidateAddressDirPath())
	log.Printf("Writing results of candidate address generation to file at '%s'.", outputPath)
	start = time.Now()
	err = addressing.WriteIPsToHexFile(outputPath, addresses)
	if err != nil {
		return err
	}
	elapsed = time.Since(start)
	generateWriteTimer.Update(elapsed)
	log.Printf("It took a total of %s to write %d addresses to file.", elapsed, len(addresses))
	outputPath = fs.GetTimedFilePath(conf.GetBloomDirPath())
	log.Printf("Writing current state of Bloom filter to file at '%s'.", outputPath)
	start = time.Now()
	err = filtering.WriteBloomFilterToFile(outputPath, bloom)
	if err != nil {
		return err
	}
	elapsed = time.Since(start)
	bloomWriteTimer.Update(elapsed)
	data.UpdateBloomFilter(bloom, outputPath)
	log.Printf("It took a total of %s to write Bloom filter to file '%s'.", elapsed, outputPath)
	return nil

}