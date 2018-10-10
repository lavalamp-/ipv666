package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"net"
	"github.com/lavalamp-/ipv666/common/addressing"
	"log"
	"github.com/rcrowley/go-metrics"
	"time"
)

var blacklistCandGenDuration = metrics.NewTimer()
var blacklistCandGenCount = metrics.NewCounter()

func init() {
	metrics.Register("blacklist_cand_gen_duration", blacklistCandGenDuration)
	metrics.Register("blacklist_cand_gen_count", blacklistCandGenCount)
}

func generateNetworkAddresses(conf *config.Configuration) (error) {
	nets, err := data.GetScanResultsNetworkRanges(conf.GetNetworkGroupDirPath())
	log.Printf("Now generating %d addresses for each of the %d blacklist network candidates.", conf.NetworkPingCount, len(nets))
	if err != nil {
		return err
	}
	var addrs []*net.IP
	start := time.Now()
	for _, networks := range nets {
		addrs = append(addrs, addressing.GenerateRandomAddressesInNetwork(networks, conf.NetworkPingCount)...)
	}
	elapsed := time.Since(start)
	blacklistCandGenDuration.Update(elapsed)
	blacklistCandGenCount.Inc(int64(len(addrs)))
	log.Printf("Successfully generated %d addresses to test for blacklist.", len(addrs))
	outputPath := getTimedFilePath(conf.GetNetworkScanTargetsDirPath())
	log.Printf("Writing results to file at path '%s'.", outputPath)
	err = addressing.WriteIPsToHexFile(outputPath, addrs)
	log.Printf("Blacklist test addresses successfully written to '%s'.", outputPath)
	if err != nil {
		return err
	}
	return nil
}
