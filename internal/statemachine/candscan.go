package statemachine

import (
	"github.com/lavalamp-/ipv666/internal/config"
	"github.com/lavalamp-/ipv666/internal/data"
	"github.com/lavalamp-/ipv666/internal/fs"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/lavalamp-/ipv666/internal/pingscan"
	"github.com/rcrowley/go-metrics"
	"github.com/spf13/viper"
	"time"
)

var liveAddrCandGauge = metrics.NewGauge()
var pingscanCandDurationTimer = metrics.NewTimer()
var pingscanCandErrorCounter = metrics.NewCounter()

func init() {
	metrics.Register("candscan.live_results.gauge", liveAddrCandGauge)
	metrics.Register("candscan.ping_scan.time", pingscanCandDurationTimer)
	metrics.Register("candscan.ping_scan.error.count", pingscanCandErrorCounter)
}

func pingScanCandidateAddresses() error {
	inputPath, err := data.GetMostRecentFilePathFromDir(config.GetCandidateAddressDirPath())
	if err != nil {
		return err
	}
	outputPath := fs.GetTimedFilePath(config.GetPingResultDirPath())
	logging.Infof(
		"Now ping-scanning IPv6 addressing found in file at path '%s'. Results will be written to '%s'.",
		inputPath,
		outputPath,
	)
	start := time.Now()
	_, err = pingscan.ScanFromConfig(inputPath, outputPath)
	elapsed := time.Since(start)
	if err != nil {
		pingscanCandErrorCounter.Inc(1)
		logging.Warnf("An error was thrown when trying to run ping-scan: %s", err)
		logging.Debugf("Ping-scan elapsed time was %s.", elapsed)
		return err
	}
	pingscanCandDurationTimer.Update(elapsed)
	liveCount, err := fs.CountLinesInFile(outputPath)
	if err != nil {
		logging.Warnf("Error when counting lines in file '%s': %e", outputPath, err)
		if viper.GetBool("ExitOnFailedMetrics") {
			return err
		}
	}
	liveAddrCandGauge.Update(int64(liveCount))
	logging.Infof("Ping-scan completed successfully in %s. Results written to file at '%s'.", elapsed, outputPath)
	return nil
}
