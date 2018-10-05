package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"log"
	"github.com/lavalamp-/ipv666/common/addresses"
	"math/rand"
	"time"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"github.com/lavalamp-/ipv666/common/shell"
	"os"
)

func zmapScanNetworkRanges(conf *config.Configuration) (error) {

	// Find the target network groups file
	netsPath, err := data.GetMostRecentFilePathFromDir(conf.GetNetworkGroupDirPath())
	if err != nil {
		return err
	}

	// Load the network groups
	log.Printf("Loading network groups")
	nets, err := addresses.GetAddressListFromHexStringsFile(netsPath)
	if err != nil {
		return err
	}

	// Generate random addresses in each network
	log.Printf("Generating %d addresses in each network range", conf.NetworkPingCount)
	rand.Seed(time.Now().UTC().UnixNano())
	file, err := ioutil.TempFile("/tmp", "addrs")
	if err != nil {
		return err
	}
	var netRanges [][]addresses.IPv6Address
	for _, net := range(nets.Addresses) {
		var netRange []addresses.IPv6Address
		for x := 0; x < conf.NetworkPingCount; x++ {
			addr := addresses.IPv6Address{net.Content}
			for x := conf.NetworkGroupingSize; x < 128; x++ {
				byteOff := (int)(x/8)
				bitOff := (uint)(x-(byteOff*8))
				byteMask := (byte)(^(rand.Intn(2) << bitOff))
				addr.Content[byteOff] |= (byte)(^byteMask)
			}
			netRange = append(netRange, addr)
			file.WriteString(fmt.Sprintf("%s\n", addr.String()))
		}
		netRanges = append(netRanges, netRange)
	}
	file.Close()

	// Scan the addresses
	inputPath, err := filepath.Abs(file.Name())
	if err != nil {
		return err
	}
	file, err = ioutil.TempFile("/tmp", "addrs-scanned")
	if err != nil {
		return err
	}
	outputPath, err := filepath.Abs(file.Name())
	if err != nil {
		return err
	}
	log.Printf(
		"Now Zmap scanning IPv6 addresses found in file at path '%s'. Results will be written to '%s'.",
		inputPath,
		outputPath,
	)
	start := time.Now()
	_, err = shell.ZmapScanFromConfig(conf, inputPath, outputPath)
	elapsed := time.Since(start)
	if err != nil {
		log.Printf("An error was thrown when trying to run zmap: %s", err)
		log.Printf("Zmap elapsed time was %s.", elapsed)
		return err
	}
	log.Printf("Zmap completed successfully in %s. Results written to file at '%s'.", elapsed, outputPath)

	// Blacklist networks with 100% response rate
	blacklistPath := getTimedFilePath(conf.GetNetworkBlacklistDirPath())
	log.Printf("Writing network blacklist to %s.", blacklistPath)
	file, err = os.OpenFile(blacklistPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	addrs, err := addresses.GetAddressListFromHexStringsFile(outputPath)
	if err != nil {
		return err
	}
	for pos, netRange := range netRanges {
		addrMiss := false
		for _, netAddr := range netRange {
			found := false
			for _, addr := range addrs.Addresses {
				if netAddr.Content == addr.Content {
					found = true
					break
				}
			}
			if found == false {
				addrMiss = true
				break
			}
		}

		// If there were no response misses blacklist this network range
		if addrMiss == false {
			file.WriteString(fmt.Sprintf("%s\n", nets.Addresses[pos].String()))
		}
	}
	file.Close()

	return nil
}
