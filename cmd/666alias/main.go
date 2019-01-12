package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/lavalamp-/ipv666/common/addressing"
	"github.com/lavalamp-/ipv666/common/blacklist"
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/fs"
	"github.com/lavalamp-/ipv666/common/setup"
	"github.com/lavalamp-/ipv666/common/shell"
	"github.com/spf13/viper"
	"log"
	"net"
)

func main() {

	var inputRange string
	var configPath string

	flag.StringVar(&inputRange,"net", "", "An IPv6 CIDR range to test as an aliased network.")
	flag.StringVar(&configPath, "config", "config.json", "Local file path to the configuration file to use.")
	flag.Parse()

	if inputRange == "" {
		log.Fatal("Please provide an input IPv6 CIDR range.")
	}

	_, cidrRange, err := net.ParseCIDR(inputRange)

	if err != nil {
		log.Fatalf("Error thrown when parsing CIDR range from '%s': %e", inputRange, err)
	} else if cidrRange == nil {
		log.Fatalf("No valid CIDR range was parsed from the value '%s'.", inputRange)
	}

	// TODO handle with Cobra

	//conf, err := config.LoadFromFile(configPath)
	//
	//if err != nil {
	//	log.Fatal("Can't proceed without loading valid configuration file.")
	//}

	err = setup.InitFilesystem()

	if err != nil {
		log.Fatal("Error thrown during filesystem initialization: ", err)
	}

	log.Printf("All systems are green. Now testing IPv6 range %s for aliased state.", cidrRange)

	ip, aliased, err := checkNetworkForAliased(cidrRange)

	if err != nil {
		log.Fatal(err)
	} else if !aliased {
		log.Fatalf("Your input range of %s does not appear to be aliased based on your current configured settings. Exiting.", inputRange)
	}

	log.Print("As the initial network appears to be aliased, we will now seek out the network length.")

	aliasedNet, err := seekAliasedNetwork(cidrRange, ip)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Aliased network found!")
	log.Println()
	log.Printf("%s", aliasedNet)

}

func seekAliasedNetwork(inputNet *net.IPNet, inputIP *net.IP) (*net.IPNet, error) {

	log.Printf("Now seeking aliased network length starting from input range of %s. Addresses that responded will be %s.", inputNet, inputIP)

	ones, _ := inputNet.Mask.Size()
	acs, err := blacklist.NewAliasCheckStates([]*net.IP{inputIP}, uint8(viper.GetInt("AliasLeftIndexStart")), uint8(ones))
	var i int
	var toReturn *net.IPNet

	if err != nil {
		log.Printf("Error thrown when creating new alias check states: %e", err)
		return nil, err
	}

	loopCount := 0
	for {
		log.Printf("Now starting loop %d.", loopCount)
		log.Print("Generating test addresses...")
		testAddrs := acs.GetTestAddresses()
		if len(testAddrs) == 0 {
			log.Printf("Somehow did not generate any test addresses in loop %d. This shouldn't happen.", loopCount)
			return nil, errors.New(fmt.Sprintf("did not generate any test addresses in loop %d", loopCount))
		}
		log.Printf("%d addresses generated for loop %d.", len(testAddrs), loopCount)
		var scanAddrs []*net.IP
		for _, testAddr := range testAddrs {
			for i = 0; i < viper.GetInt("AliasDuplicateScanCount"); i++ {
				scanAddrs = append(scanAddrs, testAddr)
			}
		}
		targetsPath := fs.GetTimedFilePath(config.GetNetworkScanTargetsDirPath())
		log.Printf("Writing %d blacklist scan addresses to file '%s'.", len(scanAddrs), targetsPath)
		err := addressing.WriteIPsToHexFile(targetsPath, scanAddrs)
		if err != nil {
			log.Printf("Error thrown when writing %d addresses to file '%s': %e", len(scanAddrs), targetsPath, err)
			return nil, err
		}
		log.Printf("Successfully wrote %d blacklist scan addresses to file '%s'.", len(scanAddrs), targetsPath)
		zmapPath := fs.GetTimedFilePath(config.GetNetworkScanResultsDirPath())
		log.Printf("Kicking off Zmap from file path '%s' to output path '%s'.", targetsPath, zmapPath)
		_, err = shell.ZmapScanFromConfig(targetsPath, zmapPath)
		if err != nil {
			log.Printf("An error was thrown when trying to run zmap: %s", err)
			return nil, err
		}
		foundAddrs, err := addressing.ReadIPsFromHexFile(zmapPath)
		if err != nil {
			log.Printf("Error thrown when reading IP addresses from file '%s': %e", zmapPath, err)
			return nil, err
		}
		log.Printf("%d addresses responded to ICMP pings.", len(foundAddrs))
		foundAddrSet := addressing.GetIPSet(foundAddrs)
		log.Printf("Updating check list with results from Zmap scan.")
		acs.Update(foundAddrSet)
		acs.PrintStates()
		if acs.GetAllFound() {
			nets, err := acs.GetAliasedNetworks()
			if err != nil {
				log.Printf("Error thrown when retrieving aliased networks from AliasCheckStates: %e", err)
				return nil, err
			} else if len(nets) == 0 {
				return nil, errors.New("no aliased network returned in call to GetAliasedNetworks (length 0)")
			}
			toReturn = nets[0]
			log.Printf("It looks like we've found the aliased network border. Aliased network is %s.", toReturn)
			break
		} else {
			log.Printf("Did not find aliased network on loop %d. Let's do this again!", loopCount)
			loopCount++
		}
	}

	log.Printf("It took a total of %d loops to identify the aliased network %s.", loopCount, toReturn)

	return toReturn, nil

}

func checkNetworkForAliased(inputNet *net.IPNet) (*net.IP, bool, error) {

	log.Printf("Now checking network range %s for aliased status.", inputNet)

	addrs := addressing.GenerateRandomAddressesInNetwork(inputNet, viper.GetInt("NetworkPingCount"))
	addrsPath := fs.GetTimedFilePath(config.GetNetworkScanTargetsDirPath())

	log.Printf("Writing %d test addresses to file at path '%s'.", len(addrs), addrsPath)

	err := addressing.WriteIPsToHexFile(addrsPath, addrs)
	if err != nil {
		log.Printf("Error thrown when writing %d addresses to file '%s': %e", len(addrs), addrsPath, err)
		return nil, false, err
	}

	log.Printf("Wrote test addresses to file at path '%s'.", addrsPath)
	zmapPath := fs.GetTimedFilePath(config.GetNetworkScanResultsDirPath())

	_, err = shell.ZmapScanFromConfig(addrsPath, zmapPath)
	if err != nil {
		log.Printf("An error was thrown when trying to run zmap: %s", err)
		return nil, false, err
	}

	foundAddrs, err := addressing.ReadIPsFromHexFile(zmapPath)
	if err != nil {
		log.Printf("Error thrown when reading IP addresses from file '%s': %e", zmapPath, err)
		return nil, false, err
	}

	threshold := (int)(float64(viper.GetInt("NetworkPingCount")) * viper.GetFloat64("NetworkBlacklistPercent"))
	log.Printf("Threshold for aliased network detection is %d (%d ping count, %f percent). %d addresses responded.", threshold, viper.GetInt("NetworkPingCount"), viper.GetFloat64("NetworkBlacklistPercent"), len(foundAddrs))

	if len(foundAddrs) >= threshold {
		log.Printf("Initial network of %s appears to be aliased.", inputNet)
		return addrs[0], true, nil
	} else {
		log.Printf("Initial network of %s does not appear to be aliased.", inputNet)
		return nil, false, nil
	}

}
