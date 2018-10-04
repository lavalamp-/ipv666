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
			// Chris
			// Generate the candidate addresses to scan from the most recent model
			err := generateCandidateAddresses(conf)
			if err != nil {
				return err
			}
		case PING_SCAN_ADDR:
			// Chris
			// Perform a Zmap scan of the candidate addresses that were generated
			err := zmapScanCandidateAddresses(conf)
			if err != nil {
				return err
			}
		case NETWORK_GROUP:
			// Marc
			// Process results of Zmap scan into a set of network ranges
		case PING_SCAN_NET:
			// Marc
			// Test each of the network ranges to see if the range responds to every IP address
		case REM_BAD_ADDR:
			// Marc
			// Remove all the addresses from the Zmap results that are in ranges that failed
			// the test in the previous step
		case UPDATE_MODEL:
			// Chris
			// Update the statistical model with the valid IPv6 results we have left over
			err := updateModelWithSuccessfulHosts(conf)
			if err != nil {
				return err
			}
		case PUSH_S3:
			// Chris
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
			// Chris
			// Remove all but the most recent files in each of the directories
		case EMIT_METRICS:
			// Chris
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
	addresses := model.GenerateMulti(conf.GenerateFirstNybble, conf.GenerateAddressCount, conf.GenerateUpdateFreq)  // TODO: filter out from blacklist
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
	// TODO read addresses from results file
	results, err := addresses.GetAddressListFromAddressesFile(resultsPath)
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
			// TODO zip file up
			log.Printf("Moving file at '%s' to S3 bucket.", zipFilePath)
			// TODO export to S3
			log.Printf("Deleting zip file at '%s'.", zipFilePath)
			// TODO delete zip file upon success
			log.Printf("Successfully moved file at '%s' to S3 with compression.", curFilePath)
		}
		log.Printf("All files in directory at '%s' processed.", curDir)
	}
	log.Printf("All %d directories successfully exported to S3.", len(allDirs))
	return nil
}
