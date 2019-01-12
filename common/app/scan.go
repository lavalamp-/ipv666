package app

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"github.com/lavalamp-/ipv666/common/fs"
	"github.com/lavalamp-/ipv666/common/statemachine"
	"github.com/rcrowley/go-metrics"
	"log"
	"time"
)

var mainLoopRunTimer = metrics.NewTimer()

func init() {
	metrics.Register("main.run.time", mainLoopRunTimer)
}

func RunScanning() {

	targetNetwork, _ := config.GetTargetNetwork()

	mostRecentNetworkString, err := data.GetMostRecentTargetNetworkString()
	if err != nil {
		log.Fatalf("Error thrown when reading most recent network string: %e", err)
	}
	if mostRecentNetworkString != targetNetwork.String() {
		if mostRecentNetworkString == "" {
			log.Printf("No prior record of a scanned network exists. Resetting state machine to scan %s appropriately.", targetNetwork)
		} else {
			log.Printf("Target network (%s) is not the most recently scanned network (%s). Resetting state machine and Bloom filter accordingly.", targetNetwork, mostRecentNetworkString)
		}
		err := statemachine.ResetStateFile(config.GetStateFilePath())
		if err != nil {
			log.Fatalf("Error thrown when resetting state file: %e", err)
		}
		_, _, err = fs.DeleteAllFilesInDirectory(config.GetBloomDirPath(), []string{})
		if err != nil {
			log.Fatalf("Error thrown when deleting Bloom directory files (path '%s'): %e", config.GetBloomDirPath(), err)
		}
		err = data.WriteMostRecentTargetNetwork(targetNetwork)
		if err != nil {
			log.Fatalf("Error thrown when writing most recent target network: %e", err)
		}
	} else {
		log.Printf("The network %s is the last network that was targeted. Picking up from where we left off.", targetNetwork)
	}

	log.Print("All systems are green. Entering state machine.")

	start := time.Now()
	err = statemachine.RunStateMachine()
	elapsed := time.Since(start)
	mainLoopRunTimer.Update(elapsed)

	//TODO push metrics

	if err != nil {
		log.Fatal(err)
	}

}