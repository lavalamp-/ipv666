package statemachine

import (
	"bufio"
	"errors"
	"github.com/lavalamp-/ipv666/internal/addressing"
	"github.com/lavalamp-/ipv666/internal/blacklist"
	"github.com/lavalamp-/ipv666/internal/config"
	"github.com/lavalamp-/ipv666/internal/data"
	"github.com/lavalamp-/ipv666/internal/fs"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/lavalamp-/ipv666/internal/pingscan"
	"github.com/rcrowley/go-metrics"
	"github.com/spf13/viper"
	"net"
	"os"
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

func newSeekPair(network *net.IPNet, addr *net.IP, count uint8) *seekPair {
	return &seekPair{
		network:	network,
		address:	addr,
		count:		count,
	}
}

func seekAliasedNetworks() error {

	logging.Infof("Starting to seek aliased networks from results of ping scan.")

	scanNets, err := data.GetScanResultsNetworkRanges()
	aliasInitialNetsCounter.Inc(int64(len(scanNets)))

	if err != nil {
		logging.Warnf("Error thrown when reading scanned networks from directory '%s': %e", config.GetNetworkGroupDirPath(), err)
		return err
	}

	seekPairs, err := checkNetworksForAliased(scanNets)
	aliasSeekPairsCounter.Inc(int64(len(seekPairs)))

	if len(seekPairs) == 0 {
		logging.Infof("None of the tested networks appeared to be aliased!")
		return nil
	}

	if err != nil {
		logging.Warnf("Error thrown when checking networks for aliased properties: %e", err)
		return err
	}

	nets, err := findAliasedNetworksFromSeekPairs(seekPairs)
	aliasAliasedNetsCount.Inc(int64(len(nets)))

	if err != nil {
		logging.Warnf("Error thrown when finding aliased networks from seek pairs: %e", err)
		return err
	}

	uniqueNets := addressing.GetUniqueNetworks(nets, viper.GetInt("LogLoopEmitFreq"))
	aliasUniqueNetsCount.Inc(int64(len(uniqueNets)))
	logging.Debugf("%d networks were found via alias seeking (%d total before de-duping).", len(uniqueNets), len(nets))

	outputPath := fs.GetTimedFilePath(config.GetAliasedNetworkDirPath())

	logging.Debugf("Writing %d aliased networks to file '%s'.", len(uniqueNets), outputPath)
	err = addressing.WriteIPv6NetworksToFile(outputPath, uniqueNets)
	if err != nil {
		logging.Warnf("Error thrown when writing to file at path '%s': %e", outputPath, err)
		return err
	}

	data.UpdateAliasedNetworks(uniqueNets, outputPath)

	logging.Infof("Successfully found %d aliased networks and wrote results to disk.", len(uniqueNets))

	return nil
}

func findAliasedNetworksFromSeekPairs(seekPairs []*seekPair) ([]*net.IPNet, error) {

	logging.Infof("Starting search for aliased networks based on %d initial starting IPs.", len(seekPairs))
	start := time.Now()

	var seekIPs []*net.IP
	for _, pair := range seekPairs {
		seekIPs = append(seekIPs, pair.address)
	}
	acs, err := blacklist.NewAliasCheckStates(seekIPs, uint8(viper.GetInt("AliasLeftIndexStart")), uint8(viper.GetInt("NetworkGroupingSize")))

	if err != nil {
		return nil, err
	}

	loopCount := 0
	var toReturn []*net.IPNet
	for {
		logging.Debugf("Now starting loop %d.", loopCount)
		err := aliasSeekLoop(acs)
		if err != nil {
			logging.Warnf("Error thrown on iteration %d of loop: %e", loopCount, err)
			return nil, err
		}
		if acs.GetAllFound() {
			toReturn, err = acs.GetAliasedNetworks()
			if err != nil {
				logging.Warnf("Error thrown when retrieving aliased networks from AliasCheckStates: %e", err)
				return nil, err
			} else if len(toReturn) == 0 {
				return nil, errors.New("no aliased network returned in call to GetAliasedNetworks (length 0)")
			}
			break
		} else {
			logging.Debugf("Did not find aliased network on loop %d. Let's do this again!", loopCount)
			loopCount++
		}
	}

	logging.Infof("It took a total of %d loops to identify all the aliased networks.", loopCount)
	aliasSeekTimer.Update(time.Since(start))
	aliasSeekLoopGauge.Update(int64(loopCount))

	return toReturn, nil

}

func aliasSeekLoop(acs *blacklist.AliasCheckStates) error {
	//TODO delete files after the function is finished?
	var i int
	start := time.Now()
	logging.Debug("Generating test addresses...")
	testAddrs := acs.GetTestAddresses()
	if len(testAddrs) == 0 {
		return errors.New("did not generate any test addresses in loop")
	}
	logging.Debugf("%d addresses generated.", len(testAddrs))
	var scanAddrs []*net.IP
	for _, testAddr := range testAddrs {
		for i = 0; i < viper.GetInt("AliasDuplicateScanCount"); i++ {
			scanAddrs = append(scanAddrs, testAddr)
		}
	}
	targetsPath := fs.GetTimedFilePath(config.GetNetworkScanTargetsDirPath())
	logging.Debugf("Writing %d blacklist scan addresses to file '%s'.", len(scanAddrs), targetsPath)
	err := addressing.WriteIPsToHexFile(targetsPath, scanAddrs)
	if err != nil {
		logging.Warnf("Error thrown when writing %d addresses to file '%s': %e", len(scanAddrs), targetsPath, err)
		return err
	}
	logging.Debugf("Successfully wrote %d blacklist scan addresses to file '%s'.", len(scanAddrs), targetsPath)
	outputPath := fs.GetTimedFilePath(config.GetNetworkScanResultsDirPath())
	logging.Debugf("Kicking off ping scan from file path '%s' to output path '%s'.", targetsPath, outputPath)
	_, err = pingscan.ScanFromConfig(targetsPath, outputPath)
	if err != nil {
		logging.Warnf("An error was thrown when running ping scan: %s", err)
		return err
	}
	foundAddrs, err := fs.ReadIPsFromHexFile(outputPath)
	if err != nil {
		logging.Warnf("Error thrown when reading IP addresses from file '%s': %e", outputPath, err)
		return err
	}
	logging.Debugf("%d addresses responded to ICMP pings.", len(foundAddrs))
	foundAddrSet := addressing.GetIPSet(foundAddrs)
	logging.Debugf("Updating check list with results from Zmap scan.")
	acs.Update(foundAddrSet)
	aliasSeekLoopTimer.Update(time.Since(start))
	return nil
}

func checkNetworksForAliased(nets []*net.IPNet) ([]*seekPair, error) {

	logging.Infof("Now testing %d networks for aliased properties.", len(nets))
	start := time.Now()

	candsPath, err := generateAliasCandidates(nets)
	if err != nil {
		return nil, err
	}

	outputPath := fs.GetTimedFilePath(config.GetNetworkScanResultsDirPath())
	logging.Debugf("Ping scanning alias candidates in file '%s'. Results will be written to '%s'.", candsPath, outputPath)

	_, err = pingscan.ScanFromConfig(candsPath, outputPath)
	if err != nil {
		return nil, err
	}
	logging.Infof("Successfully scanned alias candidates to file '%s'.", outputPath)

	foundAddrs, err := fs.ReadIPsFromHexFile(outputPath)
	if err != nil {
		return nil, err
	}

	seekPairs := getSeekPairsFromScanResults(nets, foundAddrs)
	aliasCheckTimer.Update(time.Since(start))

	return seekPairs, nil

}

func generateAliasCandidates(nets []*net.IPNet) (string, error) {

	outputPath := fs.GetTimedFilePath(config.GetNetworkScanTargetsDirPath())

	logging.Debugf("Alias checking targets will be written to file '%s'.", outputPath)

	var addrs []*net.IP
	file, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)

	for i, networks := range nets {
		if i % viper.GetInt("LogLoopEmitFreq") == 0 {
			logging.Debugf("Generating addresses for network %d out of %d.", i, len(nets))
		}
		addrs = append(addrs, addressing.GenerateRandomAddressesInNetwork(networks, viper.GetInt("NetworkPingCount"))...)
		if len(addrs) >= viper.GetInt("BlacklistFlushInterval") {
			err := flushAddressesToDisk(addrs, writer)
			if err != nil {
				return "", err
			}
			addrs = addrs[:0]
		}
	}
	if len(addrs) > 0 {
		err := flushAddressesToDisk(addrs, writer)
		if err != nil {
			return "", err
		}
	}
	writer.Flush()

	logging.Debugf("Alias candidates successfully written to file '%s'.", outputPath)
	return outputPath, nil

}

func getSeekPairsFromScanResults(nets []*net.IPNet, addrs []*net.IP) []*seekPair {

	logging.Infof("Processing %d live addresses against %d networks to find networks that appear to be aliased.", len(addrs), len(nets))

	netList := blacklist.NewNetworkBlacklist(nets)
	presenceTracker := make(map[string]*seekPair)

	for i, addr :=  range addrs {
		if i % viper.GetInt("LogLoopEmitFreq") == 0 {
			logging.Debugf("Checking address %d out of %d.", i, len(addrs))
		}
		addrNetwork := netList.GetBlacklistingNetworkFromIP(addr)
		if addrNetwork != nil {
			netString := addrNetwork.String()
			if _, ok := presenceTracker[netString]; !ok {
				presenceTracker[netString] = newSeekPair(addrNetwork, addr, 0)
			}
			presenceTracker[netString].count++
		}
	}

	logging.Debugf("Processed all %d addresses.", len(addrs))

	var toReturn []*seekPair
	threshold := (uint8)(float64(viper.GetInt("NetworkPingCount")) * viper.GetFloat64("NetworkBlacklistPercent"))

	for _, v := range presenceTracker {
		if v.count >= threshold {
			toReturn = append(toReturn, v)
		}
	}

	logging.Infof("%d (out of an initial %d) networks exhibit traits of aliased networks.", len(toReturn), len(nets))

	return toReturn
}

func flushAddressesToDisk(addrs []*net.IP, w *bufio.Writer) error {
	toWrite := addressing.GetTextLinesFromIPs(addrs)
	_, err := w.WriteString(toWrite)
	return err
}
