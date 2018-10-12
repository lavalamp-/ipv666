package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"log"
	"time"
	"net"
	"github.com/rcrowley/go-metrics"
	"github.com/lavalamp-/ipv666/common/addressing"
	"encoding/binary"	
)

var generateDurationTimer = metrics.NewTimer()
var generateBlacklistCount = metrics.NewCounter()
var generateWriteTimer = metrics.NewTimer()

func init() {
	metrics.Register("gen_duration_timer", generateDurationTimer)
	metrics.Register("gen_blacklist_count", generateBlacklistCount)
	metrics.Register("gen_write_timer", generateWriteTimer)
}

func generateCandidateAddresses(conf *config.Configuration) (error) {
	model, err := data.GetProbabilisticAddressModel(conf.GetGeneratedModelDirPath())
	if err != nil {
		return err
	}
	blacklist, err := data.GetBlacklist(conf.GetNetworkBlacklistDirPath())
	if err != nil {
		return err
	}
	log.Printf(
		"Generating a total of %d addresses based on the content of model '%s' (%d digest count). Starting nybble is %d.",
		conf.GenerateAddressCount,
		model.Name,
		model.DigestCount,
		conf.GenerateFirstNybble,
	)

	// Blacklist net.IPNet's to uint64's
	blnets := map[uint64]bool{}
	for _, ipnet:= range blacklist.Networks {
		n := binary.LittleEndian.Uint64((*ipnet).IP[:8])
		blnets[n] = true
	}

	var addresses []*net.IP
	var blacklistCount, madeCount = 0, 0
	start := time.Now()
	for len(addresses) < conf.GenerateAddressCount {
		newIP := model.GenerateSingleIP(conf.GenerateFirstNybble)
		
		n := binary.LittleEndian.Uint64((*newIP)[:8])
		_, found := blnets[n]

		// if !blacklist.IsIPBlacklisted(newIP) {
		if !found {
			addresses = append(addresses, newIP)
			madeCount++
		} else {
			blacklistCount++
		}
		if (madeCount + blacklistCount) % conf.GenerateUpdateFreq == 0 {
			log.Printf("Generated %d total addresses, %d have been valid, %d have been blacklisted.", madeCount + blacklistCount, madeCount, blacklistCount)
		}
	}
	elapsed := time.Since(start)
	generateDurationTimer.Update(elapsed)
	generateBlacklistCount.Inc(int64(blacklistCount))
	log.Printf("Took a total of %s to generate %d candidate addresses (%d blacklisted filtered out).", elapsed, conf.GenerateAddressCount, blacklistCount)
	outputPath := getTimedFilePath(conf.GetCandidateAddressDirPath())
	log.Printf("Writing results of candidate address generation to file at '%s'.", outputPath)
	start = time.Now()
	err = addressing.WriteIPsToHexFile(outputPath, addresses)
	elapsed = time.Since(start)
	generateWriteTimer.Update(elapsed)
	log.Printf("It took a total of %s to write %d addresses to file.", elapsed, len(addresses))
	if err != nil {
		return err
	}
	log.Printf("Successfully wrote %d candidate addresses to file at '%s'.", conf.GenerateAddressCount, outputPath)
	return nil
}