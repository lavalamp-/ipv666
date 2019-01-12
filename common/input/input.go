package input

import (
	"fmt"
	"github.com/lavalamp-/ipv666/common/addressing"
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/fs"
	"github.com/lavalamp-/ipv666/common/shell"
	"github.com/lavalamp-/ipv666/common/statemachine"
	"github.com/lavalamp-/ipv666/common/zrandom"
	"github.com/spf13/viper"
	"log"
	"net"
)

func PrepareFromInputFile(inputFilePath string, fileType string) error {
	// Confirm that cleaning up is ok
	if !viper.GetBool("ForceAcceptPrompts") {
		err := confirmCleanUpExisting(inputFilePath)
		if err != nil {
			return err
		}
	} else {
		log.Printf("Configured to force accept prompts. Moving forward with cleaning up prior to starting from input file '%s'.", inputFilePath)
	}

	// Load addresses from input file
	addrs, err := getIPsFromFile(inputFilePath, fileType)
	if err != nil {
		return err
	}
	// Remove IPv4 addresses
	addrs = addressing.FilterIPv4FromList(addrs)
	// Unique addresses
	addrs = removeDuplicateIPs(addrs)
	// Filter out PSLAAC addresses
	addrs = filterOutHighEntropyIPs(addrs)
	// Check that enough addresses remain
	if len(addrs) < viper.GetInt("InputMinAddresses") {
		if !viper.GetBool("ForceAcceptPrompts") {
			err := confirmTooFew(len(addrs))
			if err != nil {
				return err
			}
		} else {
			log.Printf("Configured to force accept prompts. Moving forward despite too few remaining addresses (got %d, wanted %d or more).", len(addrs), viper.GetInt("InputMinAddresses"))
		}
	}
	// Delete all existing files in all directories
	err = cleanUpWorkingDirectories()
	if err != nil {
		return err
	}
	// Write addresses to ping results file path
	writeNewAddresses(addrs)
	// Update state file to indicate that ping results should be checked for blacklist
	err = updateState()
	if err != nil {
		return err
	}
	return nil
}

func getIPsFromFile(inputFilePath string, inputFileType string) ([]*net.IP, error) {
	var toReturn []*net.IP
	var err error
	if inputFileType == "bin" {
		toReturn, err = addressing.ReadIPsFromBinaryFile(inputFilePath)
	} else {
		toReturn, err = addressing.ReadIPsFromHexFile(inputFilePath)
	}
	if err != nil {
		log.Printf("Error thrown when reading addresses from file '%s': %e", inputFilePath, err)
	} else {
		log.Printf("Successfully read %d addresses from %s file at '%s'.", len(toReturn), inputFileType, inputFilePath)
	}
	return toReturn, err
}

func updateState() error {
	err := statemachine.SetStateFile(config.GetStateFilePath(), statemachine.NETWORK_GROUP)
	if err != nil {
		log.Printf("Error thrown when attempting to update state file at path '%s': %e", config.GetStateFilePath(), err)
		return err
	}
	log.Printf("Successfully updated state file at path '%s'.", config.GetStateFilePath())
	return nil
}

func writeNewAddresses(toWrite []*net.IP) error {
	outputPath := fs.GetTimedFilePath(config.GetPingResultDirPath())
	log.Printf("Writing %d IP addresses to file at path '%s'.", len(toWrite), outputPath)
	err := addressing.WriteIPsToHexFile(outputPath, toWrite)
	if err != nil {
		log.Printf("Error thrown when writing addresses to path '%s': %e", outputPath, err)
		return err
	}
	log.Printf("Successfully wrote IP address list to '%s'.", outputPath)
	return nil
}

func confirmTooFew(count int) error {
	prompt := fmt.Sprintf("The resulting list of addresses is only %d long, and we recommend having at least %d to get good results. Continue? [y/N]", count, viper.GetInt("InputMinAddresses"))
	errMsg := fmt.Sprintf("Exiting. Please add more addresses to your input list and try again.")
	err := shell.RequireApproval(prompt, errMsg)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func confirmCleanUpExisting(inputFilePath string) error {
	prompt := fmt.Sprintf("Provided input file at path '%s'. Starting with an input file requires cleaning up all existing state from previous runs. Continue? [y/N]", inputFilePath)
	errMsg := fmt.Sprintf("Exiting. Please backup all existing state (all directories under '%s') and try again.", viper.GetString("BaseOutputDirectory"))
	err := shell.RequireApproval(prompt, errMsg)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func filterOutHighEntropyIPs(ips []*net.IP) []*net.IP {
	log.Printf("Now removing high entropy IP addresses from list of length %d (%f threshold, %d bits).", len(ips), viper.GetFloat64("InputEntropyThreshold"), viper.GetInt("InputEntropyBitLength"))
	var toReturn []*net.IP
	for i, ip := range ips {
		if i % viper.GetInt("LogLoopEmitFreq") == 0 {
			log.Printf("Processing %d out of %d for high entropy IPs.", i, len(ips))
		}
		ipBytes := ([]byte)(*ip)
		entropy := zrandom.GetEntropyOfBitsFromRight(ipBytes, viper.GetInt("InputEntropyBitLength"))
		if entropy < viper.GetFloat64("InputEntropyThreshold") {
			toReturn = append(toReturn, ip)
		}
	}
	log.Printf("Resulting list is %d long (removed %d high entropy addresses).", len(toReturn), len(ips) - len(toReturn))
	return toReturn
}

func removeDuplicateIPs(ips []*net.IP) []*net.IP {
	log.Printf("Now removing duplicates from list of IP addresses of length %d.", len(ips))
	toReturn := addressing.GetUniqueIPs(ips, viper.GetInt("LogLoopEmitFreq"))
	log.Printf("Resulting list is %d long (removed %d duplicates).", len(toReturn), len(ips) - len(toReturn))
	return toReturn
}

func cleanUpWorkingDirectories() error {
	log.Printf("Now deleting all regular files (recursively) starting in directory '%s'.", viper.GetString("BaseOutputDirectory"))
	numDeleted, numSkipped, err := fs.DeleteAllFilesInDirectory(viper.GetString("BaseOutputDirectory"), config.GetSafeFilePaths())
	if err != nil {
		log.Printf("Error thrown when deleting files under directory '%s': %e", viper.GetString("BaseOutputDirectory"), err)
		return err
	}
	log.Printf("Successfully deleted %d files (%d skipped).", numDeleted, numSkipped)
	return nil
}
