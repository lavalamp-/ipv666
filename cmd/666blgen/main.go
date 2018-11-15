package main

import (
	"flag"
	"log"
	"os"
	"github.com/lavalamp-/ipv666/common/addressing"
	"github.com/lavalamp-/ipv666/common/blacklist"
	"github.com/lavalamp-/ipv666/common/fs"
	"github.com/lavalamp-/ipv666/common/config"
	"fmt"
	"github.com/lavalamp-/ipv666/common/shell"
	"github.com/lavalamp-/ipv666/common/data"
	"net"
	"path/filepath"
)

func main() {

	var inputPath string
	var configPath string

	flag.StringVar(&inputPath,"input", "", "An input file containing IPv6 network ranges to build a blacklist from.")
	flag.StringVar(&configPath, "config", "config.json", "Local file path to the configuration file to use.")
	flag.Parse()

	if inputPath == "" {
		log.Fatal("Please provide an input file path (-input).")
	}

	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		log.Fatalf("No file found at path '%s'. Please supply a valid file path.", inputPath)
	}

	conf, err := config.LoadFromFile(configPath)

	if err != nil {
		log.Fatal("Can't proceed without loading valid configuration file.")
	}

	var newBlacklist = blacklist.NewNetworkBlacklist([]*net.IPNet{})
	fileName, err := fs.GetMostRecentFileFromDirectory(conf.GetNetworkBlacklistDirPath())

	if err != nil {
		log.Fatalf("Error thrown when reading recent files from directory '%s': %e", conf.GetNetworkBlacklistDirPath(), err)
	}

	if fileName != "" {
		filePath := filepath.Join(conf.GetNetworkBlacklistDirPath(), fileName)
		prompt := fmt.Sprintf("Existing blacklist found at path '%s'. Would you like to include its contents in the new blacklist? [y/N]", filePath)
		approved, err := shell.AskForApproval(prompt)
		if err != nil {
			log.Fatalf("Error thrown when prompting for approval: %e", err)
		}
		if approved {
			newBlacklist, err = data.GetBlacklist(conf.GetNetworkBlacklistDirPath())
			if err != nil {
				panic(err)
			}
		}
	}

	networks, err := addressing.ReadIPv6NetworksFromHexFile(inputPath)

	if err != nil {
		log.Fatalf("Error thrown when reading IPv6 networks from file '%s': %e", inputPath, err)
	}

	uniqueNetworks := addressing.GetUniqueNetworks(networks, conf.LogLoopEmitFreq)
	log.Printf("%d networks trimmed down to %d unique networks.", len(networks), len(uniqueNetworks))

	log.Printf("Adding %d network ranges to blacklist (starting with %d addresses).", len(uniqueNetworks), len(newBlacklist.Networks))

	added, skipped := newBlacklist.AddNetworks(uniqueNetworks)
	log.Printf("%d addresses were added and %d were skipped.", added, skipped)
	outputPath := fs.GetTimedFilePath(conf.GetNetworkBlacklistDirPath())

	log.Printf("Writing network blacklist with %d network ranges to file at path '%s'.", len(uniqueNetworks), outputPath)

	err = blacklist.WriteNetworkBlacklistToFile(outputPath, newBlacklist)

	if err != nil {
		log.Fatalf("Error thrown when writing blacklist to file '%s': %e", outputPath, err)
	}

	log.Printf("Successfully generated blacklist file at path '%s' using input addresses from file '%s' (list was %d long).", outputPath, inputPath, len(newBlacklist.Networks))

}
