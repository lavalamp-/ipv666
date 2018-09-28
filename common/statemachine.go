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
)


const (
	GEN_ADDRESSES	State = iota
	PING_SCAN_ADDR
	NETWORK_GROUP
	PING_SCAN_NET
	REM_BAD_ADDR
	UPDATE_MODEL
	PUSH_S3
	EMIT_METRICS
	CLEAN_UP
)

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
	if state < int(GEN_ADDRESSES) || state > int(CLEAN_UP) {
		return -1, errors.New(fmt.Sprintf("State with value %d was unexpected (expected between %d and %d, inclusive).", state, GEN_ADDRESSES, CLEAN_UP))
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
	return updateStateFile(filePath, GEN_ADDRESSES)
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
		case PUSH_S3:
			// Chris
			// Zip up all the most recent files and send them off to S3 (maintain dir structure)
		case EMIT_METRICS:
			// Chris
			// Push the metrics to wherever they need to go
		case CLEAN_UP:
			// Chris
			// Remove all but the most recent files in each of the directories
		}

		elapsed := time.Since(start)
		log.Printf("Completed state %d (took %s).", state, elapsed)

		state = (state + 1) % (CLEAN_UP + 1)
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
	inputPath, err := data.GetMostRecentCandidateFilePath(conf.GetCandidateAddressDirPath())
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
