package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"log"
	"github.com/lavalamp-/ipv666/common/addresses"
	"os"
	"fmt"
)

func getScanResultsNetworkRanges(conf *config.Configuration) (error) {

	// Find the target ping results file
	pingResultsPath, err := data.GetMostRecentFilePathFromDir(conf.GetPingResultDirPath())
	if err != nil {
		return err
	}

	// Load the ping results
	log.Printf("Loading ping scan results")
	addrs, err := addresses.GetAddressListFromHexStringsFile(pingResultsPath)
	if err != nil {
		return err
	}

	// Clear the host bits and enumerate unique networks
	outputPath := getTimedFilePath(conf.GetNetworkGroupDirPath())
	log.Printf("Writing network addresses to %s.", outputPath)
	file, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	var nets []addresses.IPv6Address
	for _, s := range(addrs.Addresses) {
		addr := addresses.IPv6Address{s.Content}
		for x := conf.NetworkGroupingSize; x < 128; x++ {
			byteOff := (int)(x/8)
			bitOff := (uint)(x-(byteOff*8))
			byteMask := (byte)(^(1 << bitOff))
			addr.Content[byteOff] &= byteMask
		}
		found := false
		for _, net := range(nets) {
			if net.Content == addr.Content {
				found = true
				break
			}
		}
		if found == false {
			nets = append(nets, addr)
		}
	}

	// Persist the networks to disk
	for _, addr := range(nets) {
		file.WriteString(fmt.Sprintf("%s\n", addr.String()));
	}

	return nil
}
