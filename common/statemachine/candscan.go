package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"github.com/lavalamp-/ipv666/common/fs"
	"github.com/lavalamp-/ipv666/common/logging"
	"github.com/lavalamp-/ipv666/common/shell"
	"github.com/rcrowley/go-metrics"
	"github.com/spf13/viper"
	"time"
)

var liveAddrCandGauge = metrics.NewGauge()
var zmapCandDurationTimer = metrics.NewTimer()
var zmapCandErrorCounter = metrics.NewCounter()

func init() {
	metrics.Register("candscan.live_results.gauge", liveAddrCandGauge)
	metrics.Register("candscan.zmap_scan.time", zmapCandDurationTimer)
	metrics.Register("candscan.zmap_scan_error.count", zmapCandErrorCounter)
}

func zmapScanCandidateAddresses() error {
	inputPath, err := data.GetMostRecentFilePathFromDir(config.GetCandidateAddressDirPath())
	if err != nil {
		return err
	}
	outputPath := fs.GetTimedFilePath(config.GetPingResultDirPath())
	logging.Infof(
		"Now Zmap scanning IPv6 addressing found in file at path '%s'. Results will be written to '%s'.",
		inputPath,
		outputPath,
	)
	start := time.Now()
	_, err = shell.ZmapScanFromConfig(inputPath, outputPath)
	elapsed := time.Since(start)
	if err != nil {
		zmapCandErrorCounter.Inc(1)
		logging.Warnf("An error was thrown when trying to run zmap: %s", err)
		logging.Debugf("Zmap elapsed time was %s.", elapsed)
		return err
	}
	zmapCandDurationTimer.Update(elapsed)
	liveCount, err := fs.CountLinesInFile(outputPath)
	if err != nil {
		logging.Warnf("Error when counting lines in file '%s': %e", outputPath, err)
		if viper.GetBool("ExitOnFailedMetrics") {
			return err
		}
	}
	liveAddrCandGauge.Update(int64(liveCount))
	logging.Infof("Zmap completed successfully in %s. Results written to file at '%s'.", elapsed, outputPath)
	return nil
}
