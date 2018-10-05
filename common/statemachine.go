package common

import (
	"github.com/lavalamp-/ipv666/common/config"
	"io/ioutil"
	"errors"
	"fmt"
	"log"
	"time"
	"github.com/lavalamp-/ipv666/common/data"
	"strconv"
	"os"
	"math/rand"
	"path/filepath"
	"github.com/lavalamp-/ipv666/common/shell"
	"github.com/lavalamp-/ipv666/common/fs"
	"github.com/lavalamp-/ipv666/common/addresses"
)


const (
	GEN_ADDRESSES	State = iota
	PING_SCAN_ADDR
	NETWORK_GROUP
	PING_SCAN_NET
	REM_BAD_ADDR
	UPDATE_MODEL
	AGGREGATE_BLACKLIST
	PUSH_S3
	CLEAN_UP
	EMIT_METRICS
)

var FIRST_STATE = GEN_ADDRESSES
var LAST_STATE = EMIT_METRICS

type State int8

func fetchStateFromFile(filePath string) (State, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return -1, err
	}
	if len(content) == 0 || len(content) > 1 {
		return -1, errors.New(fmt.Sprintf("Content of file at '%s' was of unexpected length (%d).", filePath, len(content)))
	}
	state := int(content[0])
	if state < int(FIRST_STATE) || state > int(LAST_STATE) {
		return -1, errors.New(fmt.Sprintf("State with value %d was unexpected (expected between %d and %d, inclusive).", state, FIRST_STATE, LAST_STATE))
	}
	return State(state), nil
}

func updateStateFile(filePath string, curState State) (error) {
	log.Printf("Now updating state file at path '%s' with current state of %d.", filePath, curState)
	var b []byte
	b = append(b, byte(curState))
	return ioutil.WriteFile(filePath, b, 0644)
}

func InitStateFile(filePath string) (error) {
	return updateStateFile(filePath, FIRST_STATE)
}

func RunStateMachine(conf *config.Configuration) (error) {

	log.Print("Now starting to run the state machine... Hold on to your butts.")

	state, err := fetchStateFromFile(conf.GetStateFilePath())

	if err != nil {
		return err
	}

	log.Printf("Starting at state %d.", state)

	for {

		log.Printf("Now entering state %d.", state)
		start := time.Now()

		time.Sleep(1000 * time.Millisecond)

		switch state {
		case GEN_ADDRESSES:
			// Generate the candidate addresses to scan from the most recent model
			err := generateCandidateAddresses(conf)
			if err != nil {
				return err
			}
		case PING_SCAN_ADDR:
			// Perform a Zmap scan of the candidate addresses that were generated
			err := zmapScanCandidateAddresses(conf)
			if err != nil {
				return err
			}
		case NETWORK_GROUP:
			// Process results of Zmap scan into a set of network ranges
			err := getScanResultsNetworkRanges(conf)
			if err != nil {
				return err
			}
		case PING_SCAN_NET:
			// Test each of the network ranges to see if the range responds to every IP address
			err := zmapScanNetworkRanges(conf)
			if err != nil {
				return err
			}
		case REM_BAD_ADDR:
			// Remove all the addresses from the Zmap results that are in ranges that failed
			// the test in the previous step
			err := cleanBlacklistedAddresses(conf)
			if err != nil {
				return err
			}
		case UPDATE_MODEL:
			// Update the statistical model with the valid IPv6 results we have left over
			err := updateModelWithSuccessfulHosts(conf)
			if err != nil {
				return err
			}
		case AGGREGATE_BLACKLIST:
			// Aggregate all of the blacklists into a single blacklist
		case PUSH_S3:
			// Zip up all the most recent files and send them off to S3 (maintain dir structure)
			if !conf.ExportEnabled {
				log.Printf("Exporting to S3 is disabled. Skipping export step.")
			} else {
				err := pushFilesToS3(conf)
				if err != nil {
					return err
				}
			}
		case CLEAN_UP:
			// Remove all but the most recent files in each of the directories
			if !conf.CleanUpEnabled {
				log.Printf("Clean up disabled. Skipping clean up step.")
			} else {
				err := cleanUpNonRecentFiles(conf)
				if err != nil {
					return err
				}
			}
		case EMIT_METRICS:
			// Push the metrics to wherever they need to go
		}

		elapsed := time.Since(start)
		log.Printf("Completed state %d (took %s).", state, elapsed)

		state = (state + 1) % (LAST_STATE + 1)
		err := updateStateFile(conf.GetStateFilePath(), state)
		if err != nil {
			return err
		}

	}

}

func cleanBlacklistedAddresses(conf *config.Configuration) (error) {

	// Find the blacklist file path
	blacklistPath, err := data.GetMostRecentFilePathFromDir(conf.GetNetworkBlacklistDirPath())
	if err != nil {
		return err
	}

	// Load the blacklist network addresses
	log.Printf("Loading blacklist network addresses")
	nets, err := addresses.GetAddressListFromHexStringsFile(blacklistPath)
	if err != nil {
		return err
	}

	// Find the ping results file path
	addrsPath, err := data.GetMostRecentFilePathFromDir(conf.GetPingResultDirPath())
	if err != nil {
		return err
	}

	// Load the ping results
	log.Printf("Loading ping scan result addresses")
	addrs, err := addresses.GetAddressListFromHexStringsFile(addrsPath)
	if err != nil {
		return err
	}

	// Remove addresses from blacklisted networks
	log.Printf("Removing addresses from blacklisted networks")
	var cleanAddrs []addresses.IPv6Address
	for _, addr := range(addrs.Addresses) {
		found := false
		for _, net := range(nets.Addresses) {
			match := true
			for x := 0; x < conf.NetworkGroupingSize; x++ {
				byteOff := (int)(x/8)
				bitOff := (uint)(x-(byteOff*8))
				byteMask := (byte)(1 << bitOff)
				if (addr.Content[byteOff] & byteMask) != (net.Content[byteOff] & byteMask) {
					match = false
					break
				}
			}	
			if match == true {
				found = true
				break
			}		
		}
		if found == false {
			cleanAddrs = append(cleanAddrs, addr)
		}		
	}

	// Write the clean ping response addresses to disk
	cleanPath := getTimedFilePath(conf.GetCleanPingDirPath())
	log.Printf("Writing clean addresses to %s.", cleanPath)
	file, err := os.OpenFile(cleanPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	for _, addr := range(cleanAddrs) {
		file.WriteString(fmt.Sprintf("%s\n", addr.String()))
	}
	file.Close()
	return nil	
}

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

func getTimedFilePath(baseDir string) (string) {
	curTime := strconv.FormatInt(time.Now().Unix(), 10)
	return filepath.Join(baseDir, curTime)
}

func generateCandidateAddresses(conf *config.Configuration) (error) {
	model, err := data.GetProbabilisticAddressModel(conf.GetGeneratedModelDirPath())
	if err != nil {
		return err
	}
	log.Printf(
		"Generating a total of %d addresses based on the content of model '%s' (%d digest count). Starting nybble is %d.",
		conf.GenerateAddressCount,
		model.Name,
		model.DigestCount,
		conf.GenerateFirstNybble,
	)
	start := time.Now()
	// TODO: filter out from blacklist
	addresses := model.GenerateMulti(conf.GenerateFirstNybble, conf.GenerateAddressCount, conf.GenerateUpdateFreq)
	elapsed := time.Since(start)
	log.Printf("Took a total of %s to generate %d candidate addresses", elapsed, conf.GenerateAddressCount)
	outputPath := getTimedFilePath(conf.GetCandidateAddressDirPath())
	log.Printf("Writing results of candidate address generation to file at '%s'.", outputPath)
	addresses.ToAddressesFile(outputPath, conf.GenWriteUpdateFreq)
	log.Printf("Successfully wrote %d candidate addresses to file at '%s'.", conf.GenerateAddressCount, outputPath)
	return nil
}

func zmapScanCandidateAddresses(conf *config.Configuration) (error) {
	inputPath, err := data.GetMostRecentFilePathFromDir(conf.GetCandidateAddressDirPath())
	if err != nil {
		return err
	}
	outputPath := getTimedFilePath(conf.GetPingResultDirPath())
	log.Printf(
		"Now Zmap scanning IPv6 addresses found in file at path '%s'. Results will be written to '%s'.",
		inputPath,
		outputPath,
	)
	start := time.Now()
	_, err = shell.ZmapScanFromConfig(conf, inputPath, outputPath)
	elapsed := time.Since(start)
	//  TODO do something with result
	if err != nil {
		log.Printf("An error was thrown when trying to run zmap: %s", err)
		log.Printf("Zmap elapsed time was %s.", elapsed)
		return err
	}
	log.Printf("Zmap completed successfully in %s. Results written to file at '%s'.", elapsed, outputPath)
	return nil
}

func updateModelWithSuccessfulHosts(conf *config.Configuration) (error) {
	// TODO what happens if Zmap fails silently and we keep adding the same file to our model
	resultsPath, err := data.GetMostRecentFilePathFromDir(conf.GetCleanPingDirPath())
	if err != nil {
		return err
	}
	results, err := addresses.GetAddressListFromHexStringsFile(resultsPath)
	if err != nil {
		return err
	}
	model, err := data.GetProbabilisticAddressModel(conf.GetGeneratedModelDirPath())
	if err != nil {
		return err
	}
	model.UpdateMulti(results)
	outputPath := getTimedFilePath(conf.GetGeneratedModelDirPath())
	log.Printf("Now saving updated model '%s' (%d digest count) to file at path '%s'.", model.Name, model.DigestCount, outputPath)
	model.Save(outputPath)
	log.Printf("Model '%s' was saved to file at path '%s' successfully.", model.Name, outputPath)
	return nil
}

func pushFilesToS3(conf *config.Configuration) (error) {
	// TODO break this down into multiple functions
	allDirs := conf.GetAllDirectories()
	log.Printf("Now starting to push all non-most-recent files from %d directories.", len(allDirs))
	for _, curDir := range allDirs {
		log.Printf("Processing content of directory '%s'.", curDir)
		exportFiles, err := fs.GetNonMostRecentFilesFromDirectory(curDir)
		if err != nil {
			log.Printf("Error thrown when attempting to gather files for export in directory '%s'.", curDir)
			return err
		} else if len(exportFiles) == 0 {
			log.Printf("No files found for export in directory '%s'.", curDir)
			continue
		}
		log.Printf("A total of %d files were found for export in directory '%s'.", len(exportFiles), curDir)
		for _, curFileName := range exportFiles {
			curFilePath := filepath.Join(curDir, curFileName)
			zipFilePath := fmt.Sprintf("%s.zip", curFilePath)
			log.Printf("Now exporting file at local file path '%s' to THE GREAT BEYONDDDD.", curFilePath)
			log.Printf("Zipping up file at path '%s' to file at path '%s'.", curFilePath, zipFilePath)
			err := fs.ZipFiles([]string{curFilePath}, zipFilePath)
			if err != nil {
				log.Printf("Failed to zip up file at path '%s'. Stopping export.", curFilePath)
				zipErr := os.Remove(zipFilePath)
				if zipErr != nil {
					log.Printf("Another error was thrown when trying to delete zip file at path '%s': %e", zipFilePath, err)
				}
				return err
			}
			log.Printf("Successfully created zip file at path '%s'.", zipFilePath)
			log.Printf("Moving file at '%s' to S3 bucket.", zipFilePath)
			err = data.PushFileToS3FromConfig(zipFilePath, zipFilePath, conf)
			if err != nil {
				log.Printf("Failed to move file at path '%s' to S3. Stopping export.", zipFilePath)
				zipErr := os.Remove(zipFilePath)
				if zipErr != nil {
					log.Printf("Another error was thrown when trying to delete zip file at path '%s': %e", zipFilePath, err)
				}
				return err
			}
			log.Printf("Deleting zip file at '%s'.", zipFilePath)
			err = os.Remove(zipFilePath)
			if err != nil {
				log.Printf("Error thrown when attempting to delete zip file at path '%s': %e", zipFilePath, err)
				return err
			}
			log.Printf("Successfully moved file at '%s' to S3 with compression.", curFilePath)
		}
		log.Printf("All files in directory at '%s' processed.", curDir)
	}
	log.Printf("All %d directories successfully exported to S3.", len(allDirs))
	return nil
}

func cleanUpNonRecentFiles(conf *config.Configuration) (error) {
	// TODO break this down into multiple functions
	allDirs := conf.GetAllDirectories()
	log.Printf("Now starting to delete all non-recent files from %d directories.", len(allDirs))
	for _, curDir := range allDirs {
		log.Printf("Processing content of directory '%s'.", curDir)
		exportFiles, err := fs.GetNonMostRecentFilesFromDirectory(curDir)
		if err != nil {
			log.Printf("Error thrown when attempting to gather files for deletion in directory '%s'.", curDir)
			return err
		} else if len(exportFiles) == 0 {
			log.Printf("No files found for export in directory '%s'.", curDir)
			continue
		}
		for _, curFileName := range exportFiles {
			curFilePath := filepath.Join(curDir, curFileName)
			log.Printf("Deleting file at path '%s'.", curFilePath)
			err := os.Remove(curFilePath)
			if err != nil {
				log.Printf("Error thrown when attempting to delete file at path '%s': %e", curFilePath, err)
				return err
			}
			log.Printf("Successfully deleted file at path '%s'.", curFilePath)
		}
		log.Printf("Deleted all files in directory '%s'.", curDir)
	}
	log.Printf("Successfully deleted all non-recent files from %d directories.", len(allDirs))
	return nil
}
