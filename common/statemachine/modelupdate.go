package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"log"
	"github.com/rcrowley/go-metrics"
	"time"
)

var modelUpdateDurationTimer = metrics.NewTimer()
var modelUpdateCounter = metrics.NewCounter()

func init() {
	metrics.Register("model_update_duration", modelUpdateDurationTimer)
	metrics.Register("model_update_counter", modelUpdateCounter)
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
	model.UpdateMultiIP(cleanPings, conf.ModelUpdateFreq)
	elapsed := time.Since(start)
	modelUpdateCounter.Inc(int64(len(cleanPings)))
	modelUpdateDurationTimer.Update(elapsed)
	log.Printf("Model updated with %d addresses in %s.", len(cleanPings), elapsed)
	outputPath := getTimedFilePath(conf.GetGeneratedModelDirPath())
	log.Printf("Writing new model to output path at '%s'.", outputPath)
	err = model.Save(outputPath)
	if err != nil {
		return err
	}
	log.Printf("Model successfully updated and written to '%s'.", outputPath)
	data.UpdateProbabilisticAddressModel(model, outputPath)
	return nil
}


//func updateModelWithSuccessfulHosts(conf *config.Configuration) (error) {
//	// TODO what happens if Zmap fails silently and we keep adding the same file to our model
//	resultsPath, err := data.GetMostRecentFilePathFromDir(conf.GetCleanPingDirPath())
//	if err != nil {
//		return err
//	}
//	results, err := addressing.GetAddressListFromHexStringsFile(resultsPath)
//	if err != nil {
//		return err
//	}
//	model, err := data.GetProbabilisticAddressModel(conf.GetGeneratedModelDirPath())
//	if err != nil {
//		return err
//	}
//	start := time.Now()
//	model.UpdateMulti(results)
//	elapsed := time.Since(start)
//	modelUpdateDurationTimer.Update(elapsed)
//	outputPath := getTimedFilePath(conf.GetGeneratedModelDirPath())
//	log.Printf("Now saving updated model '%s' (%d digest count) to file at path '%s'.", model.Name, model.DigestCount, outputPath)
//	model.Save(outputPath)
//	log.Printf("Model '%s' was saved to file at path '%s' successfully.", model.Name, outputPath)
//	return nil
//}
