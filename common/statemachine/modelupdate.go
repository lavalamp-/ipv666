package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"log"
	"github.com/rcrowley/go-metrics"
	"time"
	"github.com/lavalamp-/ipv666/common/fs"
)

var modelUpdateDurationTimer = metrics.NewTimer()
var modelUpdateCounter = metrics.NewCounter()

func init() {
	metrics.Register("modelupdate.update.time", modelUpdateDurationTimer)
	metrics.Register("modelupdate.update.count", modelUpdateCounter)
}

func updateModelWithSuccessfulHosts(conf *config.Configuration) (error) {
	cleanPings, err := data.GetCleanPingResults(conf.GetCleanPingDirPath())
	if err != nil {
		return err
	}
	model, err := data.GetProbabilisticAddressModel(conf.GetGeneratedModelDirPath())
	if err != nil {
		return err
	}
	log.Printf("Updating model %s with %d addresses.", model.Name, len(cleanPings))
	start := time.Now()
	model.UpdateMultiIP(cleanPings, conf.LogLoopEmitFreq, conf)
	elapsed := time.Since(start)
	modelUpdateCounter.Inc(int64(len(cleanPings)))
	modelUpdateDurationTimer.Update(elapsed)
	log.Printf("Model updated with %d addresses in %s.", len(cleanPings), elapsed)
	outputPath := fs.GetTimedFilePath(conf.GetGeneratedModelDirPath())
	log.Printf("Writing new model to output path at '%s'.", outputPath)
	err = model.Save(outputPath)
	if err != nil {
		return err
	}
	log.Printf("Model successfully updated and written to '%s'.", outputPath)
	data.UpdateProbabilisticAddressModel(model, outputPath)
	return nil
}
