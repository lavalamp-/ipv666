package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"log"
	"github.com/lavalamp-/ipv666/common/addressing"
	"net"
	"github.com/rcrowley/go-metrics"
)

var netRangesCreatedGauge = metrics.NewGauge()
var netRangesDownFromGauge = metrics.NewGauge()

func init() {
	metrics.Register("network_ranges_created", netRangesCreatedGauge)
	metrics.Register("network_ranges_down_from", netRangesDownFromGauge)
}

func getScanResultsNetworkRanges(conf *config.Configuration) (error) {
	log.Printf("Now converting ping scan for candidates into network ranges.")
	addrs, err := data.GetCandidatePingResults(conf.GetPingResultDirPath())
	if err != nil {
		return err
	}
	log.Printf("Loaded ping scan results, now converting down to addresses.")
	var nets []*net.IPNet
	for _, curAddr := range addrs.Addresses {
		byteMask := addressing.GetByteMask(conf.NetworkGroupingSize)
		nets = append(nets, &net.IPNet{
			IP:		curAddr.Content[:],
			Mask:	byteMask,
		})
	}
	nets = addressing.GetUniqueNetworks(nets)
	log.Printf("Whittled %d initial addresses down to %d network ranges with bit mask length of %d.", len(addrs.Addresses), len(nets), conf.NetworkGroupingSize)
	netRangesCreatedGauge.Update(int64(len(nets)))
	netRangesDownFromGauge.Update(int64(len(addrs.Addresses)))
	outputPath := getTimedFilePath(conf.GetNetworkGroupDirPath())
	log.Printf("Writing resulting network file to path '%s'.", outputPath)
	err = addressing.WriteIPv6NetworksToFile(outputPath, nets)
	if err != nil {
		return err
	}
	log.Printf("Resulting network file successfully written.")
	data.UpdateScanResultsNetworkRanges(nets, outputPath)
	return nil
}