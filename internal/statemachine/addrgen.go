package statemachine

import (
	"github.com/lavalamp-/ipv666/internal/addressing"
	"github.com/lavalamp-/ipv666/internal/config"
	"github.com/lavalamp-/ipv666/internal/data"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/rcrowley/go-metrics"
	"github.com/spf13/viper"
	"net"
	"time"

	"errors"
	"fmt"
	"github.com/lavalamp-/ipv666/internal/filtering"
	"github.com/lavalamp-/ipv666/internal/fs"
	bloom2 "github.com/willf/bloom"
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

func generateCandidateAddresses() error {

	// Load the statistical model, blacklist, and bloom filter

	model, err := data.GetProbabilisticAddressModel()
	if err != nil {
		return err
	}
	blacklist, err := data.GetBlacklist()
	if err != nil {
		return err
	}
	bloom, err := data.GetBloomFilter()
	if err != nil {
		return err
	}
	targetNetwork, err := config.GetTargetNetwork()
	if err != nil {
		return err
	}

	if blacklist.IsNetworkBlacklisted(targetNetwork) {
		blacklistNet := blacklist.GetBlacklistingNetworkFromNetwork(targetNetwork)
		return errors.New(fmt.Sprintf("The target network range (%s) is blaclisted (blacklisting network of %s).", targetNetwork, blacklistNet))
	}

	// Generate all of the addresses and filter out based on Bloom filter and blacklist

	logging.Infof(
		"Generating a total of %d addresses based on the content of model '%s' (%d digest count). Network range is %s.",
		viper.GetInt("GenerateAddressCount"),
		model.Name,
		model.DigestCount,
		targetNetwork,
	)
	var addresses []*net.IP
	var blacklistCount, totalBloomCount, curBloomCount, madeCount = 0, 0, 0, 0
	var bloomEmptyThreshold = int(viper.GetFloat64("BloomEmptyMultiple") * float64(viper.GetInt("GenerateAddressCount")))

	addrProcessFunc := func(toCheck *net.IP) (bool, error) {
		ipBytes := ([]byte)(*toCheck)
		var toReturn bool
		if blacklist.IsIPBlacklisted(toCheck) {
			blacklistCount++
			toReturn = true
		} else if bloom.Test(ipBytes) {
			curBloomCount++
			totalBloomCount++
			toReturn = true
		} else {
			madeCount++
			bloom.Add(ipBytes)
			toReturn = false
		}
		if (madeCount + blacklistCount + totalBloomCount) % viper.GetInt("LogLoopEmitFreq") == 0 {
			logging.Infof("Generated %d total addresses, %d have been valid, %d have been blacklisted, %d exist in Bloom filter.", madeCount + blacklistCount + totalBloomCount, madeCount, blacklistCount, totalBloomCount)
		}
		if curBloomCount >= bloomEmptyThreshold {
			logging.Infof("Bloom filter rejection rate currently exceeds threshold of %d (%d rejected). Emptying and recreating.", bloomEmptyThreshold, curBloomCount)
			bloom, err = remakeBloomFilter(addresses)
			if err != nil {
				logging.Warnf("Error thrown when remaking Bloom filter: %e", err)
				return false, err
			}
			curBloomCount = 0
			bloomEmptyCount.Inc(1)
		}
		return toReturn, nil
	}

	start := time.Now()
	targetNetwork, err = config.GetTargetNetwork()
	if err != nil {
		logging.Warnf("Error thrown when getting target network from config: %e", err)
		return err
	}
	addresses, err = model.GenerateMultiIPFromNetwork(targetNetwork, viper.GetInt("GenerateAddressCount"), addrProcessFunc)
	if err != nil {
		logging.Warnf("Error thrown when generating multiple IP addresses for network %s: %e", targetNetwork, err)
		return err
	}
	elapsed := time.Since(start)
	generateDurationTimer.Update(elapsed)
	generateBlacklistCount.Inc(int64(blacklistCount))
	generateBloomCount.Inc(int64(totalBloomCount))
	logging.Infof("Took a total of %s to generate %d candidate addresses (%d blacklisted filtered out, %d existed in Bloom filter).", elapsed, viper.GetInt("GenerateAddressCount"), blacklistCount, totalBloomCount)

	// Write addresses and Bloom filter to disk and update data manager to point to in-memory references

	outputPath := fs.GetTimedFilePath(config.GetCandidateAddressDirPath())
	logging.Debugf("Writing results of candidate address generation to file at '%s'.", outputPath)
	start = time.Now()
	err = addressing.WriteIPsToHexFile(outputPath, addresses)
	if err != nil {
		return err
	}
	elapsed = time.Since(start)
	generateWriteTimer.Update(elapsed)
	logging.Debugf("It took a total of %s to write %d addresses to file.", elapsed, len(addresses))
	outputPath = fs.GetTimedFilePath(config.GetBloomDirPath())
	logging.Debugf("Writing current state of Bloom filter to file at '%s'.", outputPath)
	start = time.Now()
	err = filtering.WriteBloomFilterToFile(outputPath, bloom)
	if err != nil {
		return err
	}
	elapsed = time.Since(start)
	bloomWriteTimer.Update(elapsed)
	data.UpdateBloomFilter(bloom, outputPath)
	logging.Debugf("It took a total of %s to write Bloom filter to file '%s'.", elapsed, outputPath)
	return nil

}

func remakeBloomFilter(existingAddrs []*net.IP) (*bloom2.BloomFilter, error) {
	logging.Debugf("Creating new Bloom filter with %d entries and %d hashes.", viper.GetInt("AddressFilterSize"), viper.GetInt("AddressFilterHashCount"))
	var filter *bloom2.BloomFilter
	if _, err := os.Stat(config.GetOutputFilePath()); !os.IsNotExist(err) {
		logging.Debugf("Output file at path '%s' exists. Creating new Bloom filter from its contents.", config.GetOutputFilePath())
		filter, err = data.LoadBloomFilterFromOutput()
		if err != nil {
			return nil, err
		}
	} else {
		logging.Debugf("No output file found at path '%s'. Starting a new Bloom filter from scratch.", config.GetOutputFilePath())
		filter = filtering.NewFromConfig()
	}
	logging.Debugf("Updating Bloom filter with %d existing addresses.", len(existingAddrs))
	for _, ip := range existingAddrs {
		ipBytes := ([]byte)(*ip)
		filter.Add(ipBytes)
	}
	logging.Debugf("Successfully created new Bloom filter and added %d existing addresses.", len(existingAddrs))
	return filter, nil
}
