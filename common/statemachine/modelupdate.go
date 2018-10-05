package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/lavalamp-/ipv666/common/data"
	"github.com/lavalamp-/ipv666/common/addresses"
	"log"
	"github.com/rcrowley/go-metrics"
	"time"
)

var modelUpdateDurationTimer = metrics.NewTimer()

func init() {
	metrics.Register("model_update_duration", modelUpdateDurationTimer)
}

func updateModelWithSuccessfulHosts(conf *config.Configuration) (error) {
	// TODO what happens if Zmap fails silently and we keep adding the same file to our model
	resultsPath, err := data.GetMostRecentFilePathFromDir(conf.GetCleanPingDirPath())
	if err != nil {
		return err
	}
	results, err := addresses.GetAddressListFromHexStringsFile(resultsPath)
	if err != nil {
		return err
	}
	model, err := data.GetProbabilisticAddressModel(conf.GetGeneratedModelDirPath())
	if err != nil {
		return err
	}
	start := time.Now()
	model.UpdateMulti(results)
	elapsed := time.Since(start)
	modelUpdateDurationTimer.Update(elapsed)
	outputPath := getTimedFilePath(conf.GetGeneratedModelDirPath())
	log.Printf("Now saving updated model '%s' (%d digest count) to file at path '%s'.", model.Name, model.DigestCount, outputPath)
	model.Save(outputPath)
	log.Printf("Model '%s' was saved to file at path '%s' successfully.", model.Name, outputPath)
	return nil
}
