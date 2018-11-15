package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"net"
	"github.com/lavalamp-/ipv666/common/addressing"
	"log"
	"github.com/rcrowley/go-metrics"
	"time"
	"github.com/lavalamp-/ipv666/common/fs"
	"os"
)

var blacklistCandGenDuration = metrics.NewTimer()
var blacklistCandGenCount = metrics.NewCounter()
var blacklistCandGenFlushDuration = metrics.NewTimer()

func init() {
	metrics.Register("blgen.cand_gen.time", blacklistCandGenDuration)
	metrics.Register("blgen.cand_gen.count", blacklistCandGenCount)
	metrics.Register("blgen.cand_file_write.time", blacklistCandGenFlushDuration)
}

func generateNetworkAddresses(conf *config.Configuration) (error) {
	nets, err := data.GetScanResultsNetworkRanges(conf.GetNetworkGroupDirPath())
	log.Printf("Now generating %d addresses for each of the %d blacklist network candidates.", conf.NetworkPingCount, len(nets))
	if err != nil {
		return err
	}
	var addrs []*net.IP
	start := time.Now()
	outputPath := fs.GetTimedFilePath(conf.GetNetworkScanTargetsDirPath())
	log.Printf("Writing results to file at path '%s'.", outputPath)
	file, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Error thrown when opening output file at path '%s': %e", outputPath, err)
		return err
	}
	defer file.Close()
	for i, networks := range nets {
		if i % conf.LogLoopEmitFreq == 0 {
			log.Printf("Generating addresses for network %d out of %d.", i, len(nets))
		}
		addrs = append(addrs, addressing.GenerateRandomAddressesInNetwork(networks, conf.NetworkPingCount)...)
		if len(addrs) >= conf.BlacklistFlushInterval {
			start := time.Now()
			toWrite := addressing.GetTextLinesFromIPs(addrs)
			_, err := file.WriteString(toWrite)
			if err != nil {
				log.Printf("Error thrown when flushing blacklist candidates to disk: %e", err)
				return err
			}
			elapsed := time.Since(start)
			blacklistCandGenFlushDuration.Update(elapsed)
			addrs = addrs[:0]
		}
	}
	if len(addrs) > 0 {
		start := time.Now()
		toWrite := addressing.GetTextLinesFromIPs(addrs)
		_, err := file.WriteString(toWrite)
		if err != nil {
			log.Printf("Error thrown when flushing blacklist candidates to disk: %e", err)
			return err
		}
		elapsed := time.Since(start)
		blacklistCandGenFlushDuration.Update(elapsed)
	}
	log.Printf("Blacklist test addresses successfully written to '%s'.", outputPath)
	elapsed := time.Since(start)
	blacklistCandGenDuration.Update(elapsed)
	blacklistCandGenCount.Inc(int64(len(addrs)))
	log.Printf("Successfully generated %d addresses to test for blacklist.", len(nets) * conf.NetworkPingCount)
	if err != nil {
		return err
	}
	return nil
}
