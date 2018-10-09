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

var liveAddrCandGauge = metrics.NewGauge()
var zmapCandDurationTimer = metrics.NewTimer()
var zmapCandErrorCounter = metrics.NewCounter()

func init() {
	metrics.Register("zmap_cand_addr_live", liveAddrCandGauge)
	metrics.Register("zmap_cand_scan_duration", zmapCandDurationTimer)
	metrics.Register("zmap_cand_scan_error_count", zmapCandErrorCounter)
}

func zmapScanCandidateAddresses(conf *config.Configuration) (error) {
	inputPath, err := data.GetMostRecentFilePathFromDir(conf.GetCandidateAddressDirPath())
	if err != nil {
		return err
	}
	outputPath := getTimedFilePath(conf.GetPingResultDirPath())
	log.Printf(
		"Now Zmap scanning IPv6 addressing found in file at path '%s'. Results will be written to '%s'.",
		inputPath,
		outputPath,
	)
	start := time.Now()
	_, err = shell.ZmapScanFromConfig(conf, inputPath, outputPath)
	elapsed := time.Since(start)
	if err != nil {
		zmapCandErrorCounter.Inc(1)
		log.Printf("An error was thrown when trying to run zmap: %s", err)
		log.Printf("Zmap elapsed time was %s.", elapsed)
		return err
	}
	zmapCandDurationTimer.Update(elapsed)
	liveCount, err := fs.CountLinesInFile(outputPath)
	if err != nil {
		log.Printf("Error when counting lines in file '%s': %e", outputPath, err)
		if conf.ExitOnFailedMetrics {
			return err
		}
	}
	liveAddrCandGauge.Update(int64(liveCount))
	log.Printf("Zmap completed successfully in %s. Results written to file at '%s'.", elapsed, outputPath)
	return nil
}
