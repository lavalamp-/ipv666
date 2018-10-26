package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"log"
	"time"
	"github.com/lavalamp-/ipv666/common/shell"
	"github.com/rcrowley/go-metrics"
	"github.com/lavalamp-/ipv666/common/fs"
)

var blCandidateCounter = metrics.NewCounter()
var blCandidateResponseCounter = metrics.NewCounter()
var zmapNetsDurationTimer = metrics.NewTimer()

func init() {
	metrics.Register("blscan.zmap.time", zmapNetsDurationTimer)
	metrics.Register("blscan.candidates.count", blCandidateCounter)
	metrics.Register("blscan.candidates_respond.count", blCandidateResponseCounter)
}

func zmapScanNetworkRanges(conf *config.Configuration) (error) {
	addrsPath, err := data.GetMostRecentFilePathFromDir(conf.GetNetworkScanTargetsDirPath())
	if err != nil {
		return err
	}
	addrsCount, err := fs.CountLinesInFile(addrsPath)
	if err != nil {
		log.Printf("Could not read lines in file '%s': %e", addrsPath, err)
		return err
	}
	blCandidateCounter.Inc(int64(addrsCount))
	log.Printf("Going to scan blacklist candidate addresses in file at path '%s' (%d addresses).", addrsPath, addrsCount)
	outputPath := fs.GetTimedFilePath(conf.GetNetworkScanResultsDirPath())
	log.Printf("Results will be written to file '%s'.", outputPath)
	start := time.Now()
	_, err = shell.ZmapScanFromConfig(conf, addrsPath, outputPath)
	elapsed := time.Since(start)
	zmapNetsDurationTimer.Update(elapsed)
	log.Printf("Zmap scan took approximately %s.", elapsed)
	if err != nil {
		log.Printf("An error was thrown when trying to run zmap: %s", err)
		return err
	}
	resultsCount, err := fs.CountLinesInFile(outputPath)
	if err != nil {
		log.Printf("Could not read lines in file '%s': %e", outputPath, err)
		return err
	}
	blCandidateResponseCounter.Inc(int64(resultsCount))
	log.Printf("%d addresses responded to the ping scan.", resultsCount)
	return nil
}
