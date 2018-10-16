package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"log"
	"github.com/lavalamp-/ipv666/common/addressing"
	"net"
	"github.com/rcrowley/go-metrics"
	"github.com/lavalamp-/ipv666/common/fs"
)

var netRangesCreatedGauge = metrics.NewGauge()
var netRangesDownFromGauge = metrics.NewGauge()

func init() {
	metrics.Register("candgroup.net_ranges.gauge", netRangesCreatedGauge)
	metrics.Register("candgroup.addrs.gauge", netRangesDownFromGauge)
}

func generateScanResultsNetworkRanges(conf *config.Configuration) (error) {
	log.Printf("Now converting ping scan for candidates into network ranges.")
	addrs, err := data.GetCandidatePingResults(conf.GetPingResultDirPath())
	if err != nil {
		return err
	}
	log.Printf("Loaded ping scan results, now converting down to networks.")
	var nets []*net.IPNet
	for _, curAddr := range addrs {
		byteMask := addressing.GetByteMask(conf.NetworkGroupingSize)
		nets = append(nets, &net.IPNet{
			IP:		*curAddr,
			Mask:	byteMask,
		})
	}
	nets = addressing.GetUniqueNetworks(nets, conf.LogLoopEmitFreq)
	log.Printf("Whittled %d initial addresses down to %d network ranges with bit mask length of %d.", len(addrs), len(nets), conf.NetworkGroupingSize)
	netRangesCreatedGauge.Update(int64(len(nets)))
	netRangesDownFromGauge.Update(int64(len(addrs)))
	outputPath := fs.GetTimedFilePath(conf.GetNetworkGroupDirPath())
	log.Printf("Writing resulting network file to path '%s'.", outputPath)
	err = addressing.WriteIPv6NetworksToFile(outputPath, nets)
	if err != nil {
		return err
	}
	log.Printf("Resulting network file successfully written.")
	data.UpdateScanResultsNetworkRanges(nets, outputPath)
	return nil
}