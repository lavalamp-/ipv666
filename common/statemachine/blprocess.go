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
	"github.com/lavalamp-/ipv666/common/fs"
	"encoding/binary"
)

var blProcessNetMembershipTimer = metrics.NewTimer()
var blNetDiscoveryCounter = metrics.NewCounter()

type blacklistTracker struct {
	Count		int
	Network		*net.IPNet
}

func init() {
	metrics.Register("blprocess.net_membership.time", blProcessNetMembershipTimer)
	metrics.Register("blprocess.new_nets.count", blNetDiscoveryCounter)
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
	// netTrackMap := make(map[string]*blacklistTracker)
	start := time.Now()
	
	// Count unique 64-bit networks
	nets := map[uint64]uint{}
	for _, resultIP := range resultIPs {
		n := binary.LittleEndian.Uint64((*resultIP)[:8])
		if _, ok := nets[n]; !ok {
    	nets[n] = 0
		}
		nets[n]++
	}

	// Identify the networks with a hit rate above the defined threshold
	blnets := []uint64{}
	threshold := (uint)(float64(conf.NetworkPingCount) * conf.NetworkBlacklistPercent)
	for n, c := range nets {
		if c >= threshold {
			blnets = append(blnets, n)
		}
	}

	elapsed := time.Since(start)
	log.Printf("Matched IP addresses and networks in %s.", elapsed)

	blProcessNetMembershipTimer.Update(elapsed)
	var blacklistNets []*net.IPNet
	// Turn the uint64 network bytes into *net.IPNet
	mask := make([]byte, 16)
	binary.LittleEndian.PutUint64(mask, uint64(0xFFFFFFFFFFFFFFFF))
	for _, n := range blnets {
		b := make([]byte, 16)
		binary.LittleEndian.PutUint64(b, uint64(n))
		blnet := net.IPNet{b, mask}
		blacklistNets = append(blacklistNets, &blnet)
	}

	log.Printf("A total of %d new blacklist networks were discovered.", len(blacklistNets))
	blNetDiscoveryCounter.Inc(int64(len(blacklistNets)))
	blacklist, err := data.GetBlacklist(conf.GetNetworkBlacklistDirPath())
	if err != nil {
		return err
	}
	log.Printf("Updating blacklist with new networks.")
	for _, curNet := range blacklistNets {
		blacklist.AddNetwork(*curNet)
	}
	outputPath := fs.GetTimedFilePath(conf.GetNetworkBlacklistDirPath())
	log.Printf("Writing new version of blacklist to file at path '%s'.", outputPath)
	err = blacklist2.WriteNetworkBlacklistToFile(outputPath, blacklist)
	if err != nil {
		return err
	}
	data.UpdateBlacklist(blacklist, outputPath)
	log.Printf("Blacklist successfully updated and written to '%s'.", outputPath)

	return nil
}
