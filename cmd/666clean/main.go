package main

import (
	"flag"
	"os"
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/addressing"
	"github.com/lavalamp-/ipv666/common/blacklist"
	"github.com/lavalamp-/ipv666/common/fs"
	"github.com/lavalamp-/ipv666/common/data"
	"log"
	"github.com/lavalamp-/ipv666/common/setup"
)

func main() {

	var inputPath string
	var configPath string
	var outputPath string
	var blacklistPath string

	flag.StringVar(&inputPath,"input", "", "An input file containing IPv6 addresses to clean via a blacklist.")
	flag.StringVar(&configPath, "config", "config.json", "Local file path to the configuration file to use.")
	flag.StringVar(&outputPath, "out", "", "The file path where the cleaned results should be written to.")
	flag.StringVar(&blacklistPath, "blacklist", "", "The local file path to the blacklist to use. If not specified, defaults to the most recent blacklist in the configured blacklist directory.")
	flag.Parse()

	// Validate input

	if inputPath == "" {
		log.Fatal("Please provide an input file path (-input).")
	}

	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		log.Fatalf("No file found at path '%s'. Please supply a valid file path.", inputPath)
	}

	if outputPath == "" {
		log.Fatal("Please provide an output file path (-out).")
	}

	if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
		log.Fatalf("File already exists at output path of '%s'. Please choose a different output path or delete the file in question.", outputPath)
	}

	if blacklistPath != "" {
		if _, err := os.Stat(blacklistPath); os.IsNotExist(err) {
			log.Fatalf("No blacklist file found at path '%s'. Please either specify a valid blacklist file path or don't specify one.", blacklistPath)
		}
	}

	// Load content from disk

	conf, err := config.LoadFromFile(configPath)

	if err != nil {
		log.Fatalf("Can't proceed without loading valid configuration file: %e", err)
	}

	err = setup.InitFilesystem(&conf)

	if err != nil {
		log.Fatal("Error thrown during filesystem initialization: ", err)
	}

	log.Printf("Loading IP addresses from file '%s'.", inputPath)

	addrs, err := addressing.ReadIPsFromHexFile(inputPath)

	if err != nil {
		log.Fatalf("Error thrown when reading input list of IP addresses at path '%s': %e", inputPath, err)
	}
	log.Printf("Successfully loaded IP addresses from '%s'.", inputPath)

	var blist *blacklist.NetworkBlacklist
	if blacklistPath != "" {
		blist, err = blacklist.ReadNetworkBlacklistFromFile(blacklistPath)
		if err != nil {
			log.Fatalf("Error thrown when reading blacklist from path '%s': %e", blacklistPath, err)
		}
	} else {
		fileName, err := fs.GetMostRecentFileFromDirectory(conf.GetNetworkBlacklistDirPath())
		if err != nil {
			log.Fatalf("Error thrown when reading recent files from directory '%s': %e", conf.GetNetworkBlacklistDirPath(), err)
		} else if fileName == "" {
			log.Fatalf("No existing blacklist found in directory %s.", conf.GetNetworkBlacklistDirPath())
		}
		blist, err = data.GetBlacklist(conf.GetNetworkBlacklistDirPath())
		if err != nil {
			log.Fatalf("Error thrown when retrieving blacklist from directory '%s': %e", conf.GetNetworkBlacklistDirPath(), err)
		}
	}

	log.Print("Successfully loaded all files from disk. Now cleaning input list.")

	// Filter out addresses

	uniqAddrs := addressing.GetUniqueIPs(addrs, conf.LogLoopEmitFreq)

	log.Printf("Whittled %d input addresses down to %d unique addresses.", len(addrs), len(uniqAddrs))

	outAddrs := blist.CleanIPList(uniqAddrs, conf.LogLoopEmitFreq)

	log.Printf("%d addresses remain after cleaning from blacklist (started with %d).", len(outAddrs), len(uniqAddrs))

	// Write results to disk

	log.Printf("Writing cleaned address list to file at path '%s'.", outputPath)

	err = addressing.WriteIPsToHexFile(outputPath, outAddrs)

	if err != nil {
		log.Fatalf("Error thrown when writing %d addresses to '%s': %e", len(outAddrs), outputPath, err)
	}

	log.Printf("Successfully wrote results to file '%s'.", outputPath)

}
