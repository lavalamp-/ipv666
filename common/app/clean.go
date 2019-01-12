package app

import (
	"github.com/lavalamp-/ipv666/common/addressing"
	"github.com/lavalamp-/ipv666/common/blacklist"
	"github.com/spf13/viper"
	"log"
)

func RunClean(inputPath string, outputPath string, blist *blacklist.NetworkBlacklist) {

	addrs, err := addressing.ReadIPsFromHexFile(inputPath)

	if err != nil {
		log.Fatalf("Error thrown when reading input list of IP addresses at path '%s': %e", inputPath, err)
	}
	log.Printf("Successfully loaded IP addresses from '%s'.", inputPath)

	uniqAddrs := addressing.GetUniqueIPs(addrs, viper.GetInt("LogLoopEmitFreq"))

	log.Printf("Whittled %d input addresses down to %d unique addresses.", len(addrs), len(uniqAddrs))

	outAddrs := blist.CleanIPList(uniqAddrs, viper.GetInt("LogLoopEmitFreq"))

	log.Printf("%d addresses remain after cleaning from blacklist (started with %d).", len(outAddrs), len(uniqAddrs))

	// Write results to disk

	log.Printf("Writing cleaned address list to file at path '%s'.", outputPath)

	err = addressing.WriteIPsToHexFile(outputPath, outAddrs)

	if err != nil {
		log.Fatalf("Error thrown when writing %d addresses to '%s': %e", len(outAddrs), outputPath, err)
	}

	log.Printf("Successfully wrote results to file '%s'.", outputPath)

}
