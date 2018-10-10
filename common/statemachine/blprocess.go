package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"github.com/lavalamp-/ipv666/common/addressing"
	"net"
	blacklist2 "github.com/lavalamp-/ipv666/common/blacklist"
	"log"
	"github.com/rcrowley/go-metrics"
	"time"
)

var blProcessNetMembershipTimer = metrics.NewTimer()
var blNetDiscoveryCounter = metrics.NewCounter()

type blacklistTracker struct {
	Count		int
	Network		*net.IPNet
}

func init() {
	metrics.Register("blacklist_process_net_membership_timer", blProcessNetMembershipTimer)
	metrics.Register("blacklist_network_discovery_count", blNetDiscoveryCounter)
}

func processBlacklistScanResults(conf *config.Configuration) (error) {
	resultsPath, err := data.GetMostRecentFilePathFromDir(conf.GetNetworkScanResultsDirPath())
	if err != nil {
		return err
	}
	log.Printf("Reading blacklist scan result IPs from file at path '%s'.", resultsPath)
	resultIPs, err := addressing.ReadIPsFromHexFile(resultsPath)
	if err != nil {
		return err
	}
	log.Printf("%d blacklist scan result IPs found in file at '%s'.", len(resultIPs), resultsPath)
	scanNetRanges, err := data.GetScanResultsNetworkRanges(conf.GetNetworkGroupDirPath())
	log.Printf("Matching %d IP addresses against %d network ranges.", len(resultIPs), len(scanNetRanges))
	netTrackMap := make(map[string]*blacklistTracker)
	start := time.Now()
	for _, resultIP := range resultIPs {
		for _, scanNetRange := range scanNetRanges {
			if scanNetRange.Contains(*resultIP) {
				if _, ok := netTrackMap[scanNetRange.String()]; ok {
					netTrackMap[scanNetRange.String()].Count++
				} else {
					netTrackMap[scanNetRange.String()] = &blacklistTracker{
						Count:		1,
						Network:	scanNetRange,
					}
				}
				break
			}
		}
	}
	elapsed := time.Since(start)
	log.Printf("Matched IP addresses and networks in %s.", elapsed)
	blProcessNetMembershipTimer.Update(elapsed)
	var blacklistNets []*net.IPNet
	for _, v := range netTrackMap {
		responsePercent := float32(v.Count) / float32(conf.NetworkPingCount)
		//TODO add histogram for tracking percentage of responses that come from nets
		if responsePercent >= conf.NetworkBlacklistPercent {
			blacklistNets = append(blacklistNets, v.Network)
		}
	}
	log.Printf("A total of %d new blacklist networks were discovered.", len(blacklistNets))
	blNetDiscoveryCounter.Inc(int64(len(blacklistNets)))
	blacklist, err := data.GetBlacklist(conf.GetNetworkBlacklistDirPath())
	if err != nil {
		return err
	}
	log.Printf("Updating blacklist with new networks.")
	for _, curNet := range blacklistNets {
		blacklist.Update(curNet)
	}
	outputPath := getTimedFilePath(conf.GetNetworkBlacklistDirPath())
	log.Printf("Writing new version of blacklist to file at path '%s'.", outputPath)
	err = blacklist2.WriteNetworkBlacklistToFile(outputPath, blacklist)
	if err != nil {
		return err
	}
	data.UpdateBlacklist(blacklist, outputPath)
	log.Printf("Blacklist successfully updated and written to '%s'.", outputPath)
	return nil
}
