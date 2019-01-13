package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"log"
	"time"
	"github.com/lavalamp-/ipv666/common/pingscan"
	"github.com/rcrowley/go-metrics"
	"github.com/lavalamp-/ipv666/common/fs"
)

var liveAddrCandGauge = metrics.NewGauge()
var pingscanCandDurationTimer = metrics.NewTimer()
var pingscanCandErrorCounter = metrics.NewCounter()

func init() {
	metrics.Register("candscan.live_results.gauge", liveAddrCandGauge)
	metrics.Register("candscan.ping_scan.time", pingscanCandDurationTimer)
	// metrics.Register("candscan.zmap_scan_error.count", zmapCandErrorCounter)
}

func pingScanCandidateAddresses(conf *config.Configuration) (error) {
	inputPath, err := data.GetMostRecentFilePathFromDir(conf.GetCandidateAddressDirPath())
	if err != nil {
		return err
	}
	outputPath := fs.GetTimedFilePath(conf.GetPingResultDirPath())
	log.Printf(
		"Now ping-scanning IPv6 addressing found in file at path '%s'. Results will be written to '%s'.",
		inputPath,
		outputPath,
	)
	start := time.Now()
	_, err = pingscan.PingScanFromConfig(conf, inputPath, outputPath)
	elapsed := time.Since(start)
	if err != nil {
		pingscanCandErrorCounter.Inc(1)
		log.Printf("An error was thrown when trying to run ping-scan: %s", err)
		log.Printf("Ping-scan elapsed time was %s.", elapsed)
		return err
	}
	pingscanCandDurationTimer.Update(elapsed)
	liveCount, err := fs.CountLinesInFile(outputPath)
	if err != nil {
		log.Printf("Error when counting lines in file '%s': %e", outputPath, err)
		if conf.ExitOnFailedMetrics {
			return err
		}
	}
	liveAddrCandGauge.Update(int64(liveCount))
	log.Printf("Ping-scan completed successfully in %s. Results written to file at '%s'.", elapsed, outputPath)
	return nil
}
