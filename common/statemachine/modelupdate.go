package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"github.com/lavalamp-/ipv666/common/fs"
	"github.com/lavalamp-/ipv666/common/logging"
	"github.com/rcrowley/go-metrics"
	"github.com/spf13/viper"
	"time"
)

var modelUpdateDurationTimer = metrics.NewTimer()
var modelUpdateCounter = metrics.NewCounter()

func init() {
	metrics.Register("modelupdate.update.time", modelUpdateDurationTimer)
	metrics.Register("modelupdate.update.count", modelUpdateCounter)
}

func updateModelWithSuccessfulHosts() error {
	cleanPings, err := data.GetCleanPingResults(config.GetCleanPingDirPath())
	if err != nil {
		return err
	}
	model, err := data.GetProbabilisticAddressModel(config.GetGeneratedModelDirPath())
	if err != nil {
		return err
	}
	logging.Infof("Updating model %s with %d addresses.", model.Name, len(cleanPings))
	start := time.Now()
	model.UpdateMultiIP(cleanPings, viper.GetInt("LogLoopEmitFreq"))
	elapsed := time.Since(start)
	modelUpdateCounter.Inc(int64(len(cleanPings)))
	modelUpdateDurationTimer.Update(elapsed)
	logging.Debugf("Model updated with %d addresses in %s.", len(cleanPings), elapsed)
	outputPath := fs.GetTimedFilePath(config.GetGeneratedModelDirPath())
	logging.Debugf("Writing new model to output path at '%s'.", outputPath)
	err = model.Save(outputPath)
	if err != nil {
		return err
	}
	logging.Debugf("Model successfully updated and written to '%s'.", outputPath)
	data.UpdateProbabilisticAddressModel(model, outputPath)
	return nil
}
