package app

import (
	"bufio"
	"fmt"
	"github.com/lavalamp-/ipv666/internal/config"
	"github.com/lavalamp-/ipv666/internal/data"
	"github.com/lavalamp-/ipv666/internal/fs"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/lavalamp-/ipv666/internal/statemachine"
	"github.com/lavalamp-/ipv666/internal/shell"
	"github.com/rcrowley/go-metrics"
	"io"
	"strconv"
	"time"
	"os"
)

var mainLoopRunTimer = metrics.NewTimer()

func init() {
	metrics.Register("main.run.time", mainLoopRunTimer)
}

func cloudSyncOptIn() error {

	// Get the cloud sync opt-in timestamp path
	path := config.GetCloudSyncOptInPath()

	// Populate the file if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
	  f, err := os.Create(path)
	  if err != nil {
	  	return err
	  }
	  f.WriteString("0\n0\n")
	  defer f.Close()
	}

	// Read the opt-in status
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	optInLine, _, err := r.ReadLine()
	if err == io.EOF {
		optInLine = []byte("0")
	} else if err != nil {
		return err
	}
	lastAskLine, _, err := r.ReadLine()
	if err == io.EOF {
		lastAskLine = []byte("0")
	} else if err != nil {
		return err
	}

	// Convert the status to ints
	optIn, err := strconv.ParseInt(string(optInLine), 10, 64)
	if err != nil {
		return err
	}
	lastAsk, err := strconv.ParseInt(string(lastAskLine), 10, 64)
	if err != nil {
		return err
	}

	// Check if already opted in
	if optIn == 1 {
		config.SetCloudSyncOptIn(true)
		return nil
	}

	// Check if the user isn't opted in to cloud sync, and it's been 7+ days since the last nag
	now := time.Now().Unix()
	if optIn == 0 && now - lastAsk > 604800 /* 7 days in seconds */ {

		// Prompt to opt-in to cloud sync
		ok, err := shell.AskForApproval("Would you like to give back to the community and contribute the cloud-sourced IPv6 dataset @ ipv6.exposed? [y/n]:")
		if err != nil {
			return err
		}
		
		// Persist the state to disk
	  ff, err := os.Create(path)
	  if err != nil {
	  	return err
	  }
	  defer ff.Close()
		if ok {
			config.SetCloudSyncOptIn(true)
			fmt.Fprintf(ff, "%d\n%d\n", 1, now)
		} else {
			config.SetCloudSyncOptIn(false)
			fmt.Fprintf(ff, "%d\n%d\n", 0, now)
		}
	}
	return nil
}


// TODO add functionality for writing results in hex format

func RunDiscovery() {

	cloudSyncOptIn()

	targetNetwork, _ := config.GetTargetNetwork()

	mostRecentNetworkString, err := data.GetMostRecentTargetNetworkString()
	if err != nil {
		logging.ErrorStringFf("Error thrown when reading most recent network string: %e", err)
	}
	if mostRecentNetworkString != targetNetwork.String() {
		if mostRecentNetworkString == "" {
			logging.Infof("No prior record of a scanned network exists. Resetting state machine to scan %s appropriately.", targetNetwork)
		} else {
			logging.Infof("Target network (%s) is not the most recently scanned network (%s). Resetting state machine and Bloom filter accordingly.", targetNetwork, mostRecentNetworkString)
		}
		err := statemachine.ResetStateFile(config.GetStateFilePath())
		if err != nil {
			logging.ErrorStringFf("Error thrown when resetting state file: %e", err)
		}
		_, _, err = fs.DeleteAllFilesInDirectory(config.GetBloomDirPath(), []string{})
		if err != nil {
			logging.ErrorStringFf("Error thrown when deleting Bloom directory files (path '%s'): %e", config.GetBloomDirPath(), err)
		}
		err = data.WriteMostRecentTargetNetwork(targetNetwork)
		if err != nil {
			logging.ErrorStringFf("Error thrown when writing most recent target network: %e", err)
		}
	} else {
		logging.Infof("The network %s is the last network that was targeted. Picking up from where we left off.", targetNetwork)
	}

	logging.Info("All systems are green. Entering state machine.")

	start := time.Now()
	err = statemachine.RunStateMachine()
	elapsed := time.Since(start)
	mainLoopRunTimer.Update(elapsed)

	//TODO push metrics

	if err != nil {
		logging.ErrorF(err)
	}

}