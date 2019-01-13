package statemachine

import (
	"github.com/lavalamp-/ipv666/common/addressing"
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"github.com/lavalamp-/ipv666/common/fs"
	"github.com/lavalamp-/ipv666/common/logging"
	"github.com/rcrowley/go-metrics"
	"github.com/spf13/viper"
	"net"
)

var netRangesCreatedGauge = metrics.NewGauge()
var netRangesDownFromGauge = metrics.NewGauge()

func init() {
	metrics.Register("candgroup.net_ranges.gauge", netRangesCreatedGauge)
	metrics.Register("candgroup.addrs.gauge", netRangesDownFromGauge)
}

func generateScanResultsNetworkRanges() error {
	logging.Infof("Now converting ping scan for candidates into network ranges.")
	addrs, err := data.GetCandidatePingResults()
	if err != nil {
		return err
	}
	logging.Debugf("Loaded ping scan results, now converting down to networks.")
	var nets []*net.IPNet
	for _, curAddr := range addrs {
		newNet, err := addressing.GetIPv6NetworkFromBytes(*curAddr, uint8(viper.GetInt("NetworkGroupingSize")))
		if err != nil {
			return err
		}
		nets = append(nets, newNet)
	}
	nets = addressing.GetUniqueNetworks(nets, viper.GetInt("LogLoopEmitFreq"))
	logging.Debugf("Whittled %d initial addresses down to %d network ranges with bit mask length of %d.", len(addrs), len(nets), viper.GetInt("NetworkGroupingSize"))
	netRangesCreatedGauge.Update(int64(len(nets)))
	netRangesDownFromGauge.Update(int64(len(addrs)))
	outputPath := fs.GetTimedFilePath(config.GetNetworkGroupDirPath())
	logging.Debugf("Writing resulting network file to path '%s'.", outputPath)
	err = addressing.WriteIPv6NetworksToFile(outputPath, nets)
	if err != nil {
		return err
	}
	logging.Debugf("Resulting network file successfully written.")
	data.UpdateScanResultsNetworkRanges(nets, outputPath)
	return nil
}