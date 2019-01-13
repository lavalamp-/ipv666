package app

import (
	"github.com/lavalamp-/ipv666/internal/config"
	"github.com/lavalamp-/ipv666/internal/data"
	"github.com/lavalamp-/ipv666/internal/fs"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/lavalamp-/ipv666/internal/statemachine"
	"github.com/rcrowley/go-metrics"
	"time"
)

var mainLoopRunTimer = metrics.NewTimer()

func init() {
	metrics.Register("main.run.time", mainLoopRunTimer)
}

// TODO add functionality for writing results in hex format

func RunDiscovery() {

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