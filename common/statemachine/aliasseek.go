package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"log"
	"net"
	"github.com/lavalamp-/ipv666/common/fs"
	"os"
	"github.com/lavalamp-/ipv666/common/addressing"
	"github.com/lavalamp-/ipv666/common/shell"
	"github.com/lavalamp-/ipv666/common/data"
	"github.com/lavalamp-/ipv666/common/blacklist"
	"errors"
	"github.com/rcrowley/go-metrics"
	"time"
)

var aliasCheckTimer = metrics.NewTimer()
var aliasSeekTimer = metrics.NewTimer()
var aliasSeekLoopGauge = metrics.NewGauge()
var aliasSeekLoopTimer = metrics.NewTimer()
var aliasInitialNetsCounter = metrics.NewCounter()
var aliasSeekPairsCounter = metrics.NewCounter()
var aliasAliasedNetsCount = metrics.NewCounter()
var aliasUniqueNetsCount = metrics.NewCounter()

func init() {
	metrics.Register("aliasseek.aliascheck.time", aliasCheckTimer)
	metrics.Register("aliasseek.aliasseek.time", aliasSeekTimer)
	metrics.Register("aliasseek.aliaseek.loop.gauge", aliasSeekLoopGauge)
	metrics.Register("aliasseek.aliasseek.loop.time", aliasSeekLoopTimer)
	metrics.Register("aliasseek.scannets.count", aliasInitialNetsCounter)
	metrics.Register("aliasseek.seekpairs.count", aliasSeekPairsCounter)
	metrics.Register("aliasseek.foundnets.count", aliasAliasedNetsCount)
	metrics.Register("aliasseek.uniquefoundnets.count", aliasUniqueNetsCount)
}

type seekPair struct {
	network		*net.IPNet
	address		*net.IP
	count		uint8
}

func newSeekPair(network *net.IPNet, addr *net.IP, count uint8) (*seekPair) {
	return &seekPair{
		network:	network,
		address:	addr,
		count:		count,
	}
}

func seekAliasedNetworks(conf *config.Configuration) (error) {

	log.Print("Starting to seek aliased networks from results of Zmap scan.")

	scanNets, err := data.GetScanResultsNetworkRanges(conf.GetNetworkGroupDirPath())
	aliasInitialNetsCounter.Inc(int64(len(scanNets)))

	if err != nil {
		log.Printf("Error thrown when reading scanned networks from directory '%s': %e", conf.GetNetworkGroupDirPath(), err)
		return err
	}

	seekPairs, err := checkNetworksForAliased(scanNets, conf)
	aliasSeekPairsCounter.Inc(int64(len(seekPairs)))

	if err != nil {
		log.Printf("Error thrown when checking networks for aliased properties: %e", err)
		return err
	}

	nets, err := findAliasedNetworksFromSeekPairs(seekPairs, conf)
	aliasAliasedNetsCount.Inc(int64(len(nets)))

	if err != nil {
		log.Printf("Error thrown when finding aliased networks from seek pairs: %e", err)
		return err
	}

	uniqueNets := addressing.GetUniqueNetworks(nets, conf.LogLoopEmitFreq)
	aliasUniqueNetsCount.Inc(int64(len(uniqueNets)))
	log.Printf("%d networks were found via alias seeking (%d total before de-duping).", len(uniqueNets), len(nets))

	outputPath := fs.GetTimedFilePath(conf.GetAliasedNetworkDirPath())

	log.Printf("Writing %d aliased networks to file '%s'.", len(uniqueNets), outputPath)
	err = addressing.WriteIPv6NetworksToFile(outputPath, uniqueNets)
	if err != nil {
		log.Printf("Error thrown when writing to file at path '%s': %e", outputPath, err)
		return err
	}

	log.Printf("Successfully found %d aliased networks and wrote results to disk.", len(uniqueNets))

	return nil
}

func findAliasedNetworksFromSeekPairs(seekPairs []*seekPair, conf *config.Configuration) ([]*net.IPNet, error) {

	log.Printf("Starting search for aliased networks based on %d initial starting IPs.", len(seekPairs))
	start := time.Now()

	var seekIPs []*net.IP
	for _, pair := range seekPairs {
		seekIPs = append(seekIPs, pair.address)
	}
	acs, err := blacklist.NewAliasCheckStates(seekIPs, conf.AliasLeftIndexStart, conf.NetworkGroupingSize)

	if err != nil {
		return nil, err
	}

	loopCount := 0
	var toReturn []*net.IPNet
	for {
		log.Printf("Now starting loop %d.", loopCount)
		err := aliasSeekLoop(acs, conf)
		if err != nil {
			log.Printf("Error thrown on iteration %d of loop: %e", loopCount, err)
			return nil, err
		}
		if acs.GetAllFound() {
			toReturn, err := acs.GetAliasedNetworks()
			if err != nil {
				log.Printf("Error thrown when retrieving aliased networks from AliasCheckStates: %e", err)
				return nil, err
			} else if len(toReturn) == 0 {
				return nil, errors.New("no aliased network returned in call to GetAliasedNetworks (length 0)")
			}
			log.Printf("It looks like we've found the aliased network border. Aliased network is %s.", toReturn)
			break
		} else {
			log.Printf("Did not find aliased network on loop %d. Let's do this again!", loopCount)
			loopCount++
		}
	}

	log.Printf("It took a total of %d loops to identify all the aliased networks.", loopCount)
	aliasSeekTimer.Update(time.Since(start))
	aliasSeekLoopGauge.Update(int64(loopCount))

	return toReturn, nil

}

func aliasSeekLoop(acs *blacklist.AliasCheckStates, conf *config.Configuration) (error) {
	//TODO delete files after the function is finished?
	var i uint8
	start := time.Now()
	log.Print("Generating test addresses...")
	testAddrs := acs.GetTestAddresses()
	if len(testAddrs) == 0 {
		return errors.New("did not generate any test addresses in loop %d")
	}
	log.Printf("%d addresses generated.", len(testAddrs))
	var scanAddrs []*net.IP
	for _, testAddr := range testAddrs {
		for i = 0; i < conf.AliasDuplicateScanCount; i++ {
			scanAddrs = append(scanAddrs, testAddr)
		}
	}
	toWrite := addressing.GetTextLinesFromIPs(scanAddrs)
	targetsPath := fs.GetTimedFilePath(conf.GetNetworkScanTargetsDirPath())
	log.Printf("Writing %d blacklist scan addresses to file '%s'.", len(scanAddrs), targetsPath)
	targetsFile, err := os.OpenFile(targetsPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Error thrown when opening output file at path '%s': %e", targetsPath, err)
		return err
	}
	_, err = targetsFile.WriteString(toWrite)
	if err != nil {
		log.Printf("Error thrown when flushing blacklist candidates to disk: %e", err)
		return err
	}
	targetsFile.Close()
	log.Printf("Successfully wrote %d blacklist scan addresses to file '%s'.", len(scanAddrs), targetsPath)
	zmapPath := fs.GetTimedFilePath(conf.GetNetworkScanResultsDirPath())
	log.Printf("Kicking off Zmap from file path '%s' to output path '%s'.", targetsPath, zmapPath)
	_, err = shell.ZmapScanFromConfig(conf, targetsPath, zmapPath)
	if err != nil {
		log.Printf("An error was thrown when trying to run zmap: %s", err)
		return err
	}
	foundAddrs, err := addressing.ReadIPsFromHexFile(zmapPath)
	if err != nil {
		log.Printf("Error thrown when reading IP addresses from file '%s': %e", zmapPath, err)
		return err
	}
	log.Printf("%d addresses responded to ICMP pings.", len(foundAddrs))
	foundAddrSet := addressing.GetIPSet(foundAddrs)
	log.Printf("Updating check list with results from Zmap scan.")
	acs.Update(foundAddrSet)
	aliasSeekLoopTimer.Update(time.Since(start))
	return nil
}

func checkNetworksForAliased(nets []*net.IPNet, conf *config.Configuration) ([]*seekPair, error) {

	log.Printf("Now testing %d networks for aliased properties.", len(nets))
	start := time.Now()

	candsPath, err := generateAliasCandidates(nets, conf)
	if err != nil {
		return nil, err
	}

	zmapPath := fs.GetTimedFilePath(conf.GetNetworkScanResultsDirPath())
	log.Printf("Zmap scanning alias candidates in file '%s'. Results will be written to '%s'.", candsPath, zmapPath)

	_, err = shell.ZmapScanFromConfig(conf, candsPath, zmapPath)
	if err != nil {
		return nil, err
	}
	log.Printf("Successfull scanned alias candidates to file '%s'.", zmapPath)

	foundAddrs, err := addressing.ReadIPsFromHexFile(zmapPath)
	if err != nil {
		return nil, err
	}

	seekPairs := getSeekPairsFromScanResults(nets, foundAddrs, conf)
	aliasCheckTimer.Update(time.Since(start))

	return seekPairs, nil

}

func generateAliasCandidates(nets []*net.IPNet, conf *config.Configuration) (string, error) {

	outputPath := fs.GetTimedFilePath(conf.GetNetworkScanTargetsDirPath())

	log.Printf("Alias checking targets will be written to file '%s'.", outputPath)

	var addrs []*net.IP
	file, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	for i, networks := range nets {
		if i % conf.LogLoopEmitFreq == 0 {
			log.Printf("Generating addresses for network %d out of %d.", i, len(nets))
		}
		addrs = append(addrs, addressing.GenerateRandomAddressesInNetwork(networks, conf.NetworkPingCount)...)
		if len(addrs) >= conf.BlacklistFlushInterval {
			err := flushAddressesToDisk(addrs, file)
			if err != nil {
				return "", err
			}
			addrs = addrs[:0]
		}
	}
	if len(addrs) > 0 {
		err := flushAddressesToDisk(addrs, file)
		if err != nil {
			return "", err
		}
	}

	log.Printf("Alias candidates successfully written to file '%s'.", outputPath)
	return outputPath, nil

}

func getSeekPairsFromScanResults(nets []*net.IPNet, addrs []*net.IP, conf *config.Configuration) ([]*seekPair) {

	log.Printf("Processing %d live addresses against %d networks to find networks that appear to be aliased.", len(addrs), len(nets))

	netList := blacklist.NewNetworkBlacklist(nets)
	presenceTracker := make(map[string]*seekPair)

	for i, addr :=  range addrs {
		if i % conf.LogLoopEmitFreq == 0 {
			log.Printf("Checking address %d out of %d.", i, len(addrs))
		}
		addrNetwork := netList.GetBlacklistingNetwork(addr)
		if addrNetwork != nil {
			netString := addrNetwork.String()
			if _, ok := presenceTracker[netString]; !ok {
				presenceTracker[netString] = newSeekPair(addrNetwork, addr, 0)
			}
			presenceTracker[netString].count++
		}
	}

	log.Printf("Processed all %d addresses.", len(addrs))

	var toReturn []*seekPair
	threshold := (uint8)(float64(conf.NetworkPingCount) * conf.NetworkBlacklistPercent)

	for _, v := range presenceTracker {
		if v.count >= threshold {
			toReturn = append(toReturn, v)
		}
	}

	log.Printf("%d (out of an initial %d) networks exhibit traits of aliased networks.", len(toReturn), len(nets))

	return toReturn
}

func flushAddressesToDisk(addrs []*net.IP, f *os.File) (error) {
	toWrite := addressing.GetTextLinesFromIPs(addrs)
	_, err := f.WriteString(toWrite)
	return err
}
