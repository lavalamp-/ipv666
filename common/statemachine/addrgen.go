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
	"github.com/willf/bloom"
	"os"
)

var generateDurationTimer = metrics.NewTimer()
var generateBlacklistCount = metrics.NewCounter()
var generateBloomCount = metrics.NewCounter()
var generateWriteTimer = metrics.NewTimer()
var bloomWriteTimer = metrics.NewTimer()
var bloomEmptyCount = metrics.NewCounter()

func init() {
	metrics.Register("addrgen.generate_duration.time", generateDurationTimer)
	metrics.Register("addrgen.generate_blacklist.count", generateBlacklistCount)
	metrics.Register("addrgen.generate_bloom.count", generateBloomCount)
	metrics.Register("addrgen.candidate_write.time", generateWriteTimer)
	metrics.Register("addrgen.bloom_write.time", bloomWriteTimer)
	metrics.Register("addrgen.bloom_empty.count", bloomEmptyCount)
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
	var blacklistCount, totalBloomCount, curBloomCount, madeCount = 0, 0, 0, 0
	var bloomEmptyThreshold = int(conf.BloomEmptyMultiple * float64(conf.GenerateAddressCount))
	start := time.Now()
	for len(addresses) < conf.GenerateAddressCount {
		newIP := model.GenerateSingleIP(conf.GenerateFirstNybble)
		ipBytes := ([]byte)(*newIP)
		if blacklist.IsIPBlacklisted(newIP) {
			blacklistCount++
		} else if bloom.Test(ipBytes) {
			curBloomCount++
			totalBloomCount++
		} else {
			madeCount++
			addresses = append(addresses, newIP)
			bloom.Add(ipBytes)
		}
		if (madeCount + blacklistCount + totalBloomCount) % conf.GenerateUpdateFreq == 0 {
			log.Printf("Generated %d total addresses, %d have been valid, %d have been blacklisted, %d exist in Bloom filter.", madeCount + blacklistCount, madeCount, blacklistCount, totalBloomCount)
		}
		if curBloomCount >= bloomEmptyThreshold {
			log.Printf("Bloom filter rejection rate currently exceeds threshold of %d (%d rejected). Emptying and recreating.", bloomEmptyThreshold, curBloomCount)
			bloom, err = remakeBloomFilter(conf, addresses)
			if err != nil {
				log.Printf("Error thrown when remaking Bloom filter: %e", err)
				return err
			}
			curBloomCount = 0
			bloomEmptyCount.Inc(1)
		}
	}
	elapsed := time.Since(start)
	generateDurationTimer.Update(elapsed)
	generateBlacklistCount.Inc(int64(blacklistCount))
	generateBloomCount.Inc(int64(totalBloomCount))
	log.Printf("Took a total of %s to generate %d candidate addresses (%d blacklisted filtered out, %d existed in Bloom filter).", elapsed, conf.GenerateAddressCount, blacklistCount, totalBloomCount)

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

func remakeBloomFilter(conf *config.Configuration, existingAddrs []*net.IP) (*bloom.BloomFilter, error) {
	log.Printf("Creating new Bloom filter with %d entries and %d hashes.", conf.AddressFilterSize, conf.AddressFilterHashCount)
	var filter *bloom.BloomFilter
	if _, err := os.Stat(conf.GetOutputFilePath()); !os.IsNotExist(err) {
		log.Printf("Output file at path '%s' exists. Creating new Bloom filter from its contents.", conf.GetOutputFilePath())
		filter, err = data.LoadBloomFilterFromOutput(conf)
		if err != nil {
			return nil, err
		}
	} else {
		log.Printf("No output file found at path '%s'. Starting a new Bloom filter from scratch.", conf.GetOutputFilePath())
		filter = filtering.NewFromConfig(conf)
	}
	log.Printf("Updating Bloom filter with %d existing addresses.", len(existingAddrs))
	for _, ip := range existingAddrs {
		ipBytes := ([]byte)(*ip)
		filter.Add(ipBytes)
	}
	log.Printf("Successfully created new Bloom filter and added %d existing addresses.", len(existingAddrs))
	return filter, nil
}
