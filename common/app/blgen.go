package app

import (
	"github.com/lavalamp-/ipv666/common/addressing"
	"github.com/lavalamp-/ipv666/common/blacklist"
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"github.com/lavalamp-/ipv666/common/fs"
	"github.com/lavalamp-/ipv666/common/shell"
	"github.com/spf13/viper"
	"log"
	"net"
)

func RunBlgen(inputPath string) {

	var newBlacklist *blacklist.NetworkBlacklist

	approved, err := shell.AskForApproval("Would you like to add to the existing blacklist (if not, a new one will be created)? [y/N]")

	if err != nil {
		log.Fatal(err)
	}

	if approved {
		log.Printf("Loading existing blacklist...")
		newBlacklist, err = data.GetBlacklist(config.GetNetworkBlacklistDirPath())
		if err != nil {
			log.Fatal(err)
		}
	} else {
		newBlacklist = blacklist.NewNetworkBlacklist([]*net.IPNet{})
	}

	networks, err := addressing.ReadIPv6NetworksFromHexFile(inputPath)

	if err != nil {
		log.Fatalf("Error thrown when reading IPv6 networks from file '%s': %e", inputPath, err)
	}

	uniqueNetworks := addressing.GetUniqueNetworks(networks, viper.GetInt("LogLoopEmitFreq"))
	log.Printf("%d networks trimmed down to %d unique networks.", len(networks), len(uniqueNetworks))

	log.Printf("Adding %d network ranges to blacklist (starting with %d addresses).", len(uniqueNetworks), newBlacklist.GetCount())

	added, skipped := newBlacklist.AddNetworks(uniqueNetworks)
	log.Printf("%d addresses were added and %d were skipped.", added, skipped)

	startCount := newBlacklist.GetCount()
	newBlacklist.Clean(viper.GetInt("LogLoopEmitFreq"))
	log.Printf("Cleaned up duplicated networks from blacklist. Down to %d networks (from %d).", newBlacklist.GetCount(), startCount)

	outputPath := fs.GetTimedFilePath(config.GetNetworkBlacklistDirPath())

	log.Printf("Writing network blacklist with %d network ranges to file at path '%s'.", newBlacklist.GetCount(), outputPath)

	err = blacklist.WriteNetworkBlacklistToFile(outputPath, newBlacklist)

	if err != nil {
		log.Fatalf("Error thrown when writing blacklist to file '%s': %e", outputPath, err)
	}

	log.Printf("Successfully generated blacklist file at path '%s' using input addresses from file '%s' (list was %d long).", outputPath, inputPath, newBlacklist.GetCount())

}
