package statemachine

import (
	"log"
	"time"
	"io/ioutil"
	"errors"
	"fmt"
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/rcrowley/go-metrics"
)

const (
	GEN_ADDRESSES	State = iota
	PING_SCAN_ADDR
	NETWORK_GROUP
	GEN_NET_ADDRS
	PING_SCAN_NET
	PROCESS_BL_CHECK_RESULTS
	REM_BAD_ADDR
	UPDATE_MODEL
	UPDATE_ADDR_FILE
	PUSH_S3
	CLEAN_UP
	EMIT_METRICS
)

var FIRST_STATE = GEN_ADDRESSES
var LAST_STATE = EMIT_METRICS

type State int8

var stateLoopTimers = make(map[string]metrics.Timer)

func getStateLoopTimer(state State, conf *config.Configuration) (metrics.Timer, bool) {
	key := fmt.Sprintf("%s.timer.%d", conf.MetricsStateLoopPrefix, state)
	if _, ok := stateLoopTimers[key]; !ok {
		timer := metrics.NewTimer()
		metrics.Register(key, timer)
		stateLoopTimers[key] = timer
	}
	val, found := stateLoopTimers[key]
	return val, found
}

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

func SetStateFile(filePath string, curState State) (error) {
	log.Printf("Now updating state file at path '%s' with current state of %d.", filePath, curState)
	var b []byte
	b = append(b, byte(curState))
	return ioutil.WriteFile(filePath, b, 0644)
}

func InitStateFile(filePath string) (error) {
	return SetStateFile(filePath, FIRST_STATE)
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

		switch state {
		case GEN_ADDRESSES:
			// Generate the candidate addressing to scan from the most recent model
			err := generateCandidateAddresses(conf) // Looking gr8
			if err != nil {
				return err
			}
		case PING_SCAN_ADDR:
			// Perform a Zmap scan of the candidate addressing that were generated
			err := zmapScanCandidateAddresses(conf) // Looking gr8
			if err != nil {
				return err
			}
		case NETWORK_GROUP:
			// Process results of Zmap scan into a set of network ranges
			err := generateScanResultsNetworkRanges(conf) // Looking gr8
			if err != nil {
				return err
			}
		case GEN_NET_ADDRS:
			// Generate addresses from the network ranges assembled in the previous step
			err := generateNetworkAddresses(conf) // Looking gr8
			if err != nil {
				return err
			}
		case PING_SCAN_NET:
			// Test each of the network ranges to see if the range responds to every IP address
			err := zmapScanNetworkRanges(conf) // Looking gr8
			if err != nil {
				return err
			}
		case PROCESS_BL_CHECK_RESULTS:
			// Process the results of the network range ping scans and update the blacklist
			err := processBlacklistScanResults(conf) // Looking gr8
			if err != nil {
				return err
			}
		case REM_BAD_ADDR:
			// Remove all the addressing from the Zmap results that are in ranges that failed
			// the test in the previous step
			err := cleanBlacklistedAddresses(conf) // Looking gr8
			if err != nil {
				return err
			}
		case UPDATE_MODEL:
			// Update the statistical model with the valid IPv6 results we have left over
			err := updateModelWithSuccessfulHosts(conf) // Looking gr8
			if err != nil {
				return err
			}
		case UPDATE_ADDR_FILE:
			// Update the cumulative addresses file
			err := updateAddressFile(conf)
			if err != nil {
				return err
			}
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
			return nil
			// Push the metrics to wherever they need to go
		}

		elapsed := time.Since(start)
		log.Printf("Completed state %d (took %s).", state, elapsed)

		timer, found := getStateLoopTimer(state, conf)
		if !found {
			log.Printf("Unable to find state loop timer for state %d.", state)
			if conf.ExitOnFailedMetrics {
				return errors.New(fmt.Sprintf("Unable to find state loop timer for state %d.", state))
			}
		}
		timer.Update(elapsed)

		state = (state + 1) % (LAST_STATE + 1)
		err = SetStateFile(conf.GetStateFilePath(), state)
		if err != nil {
			return err
		}

	}

}
