package app

import (
	"github.com/lavalamp-/ipv666/internal/fs"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/lavalamp-/ipv666/internal/modeling"
	"github.com/spf13/viper"
)

func RunModelgen(inputPath string, outputPath string, density float64) {

	logging.Infof("Reading source addresses from file at path '%s'.", inputPath)

	var err error
	addrs, err := fs.ReadIPsFromFile(inputPath)

	if err != nil {
		logging.ErrorF(err)
	}
	logging.Debugf("Successfully read %d addresses from file '%s'.", len(addrs), inputPath)

	modelSizePercent := viper.GetFloat64("ModelDefaultSizePercent")
	modelPickPercent := viper.GetFloat64("ModelDefaultPickThreshold")

	modelSize := int(modelSizePercent * float64(len(addrs)))
	pickCount := int(modelPickPercent * float64(modelSize))

	modelSize = 500000
	pickSize := 10000
	pickCount = 2500

	//if (modelSize / pickCount) * pickSize > len(addrs) {
	//	logging.ErrorStringFf("With a size percentage of %f and a pick percentage of %f, there are not enough candidates generate the initial seeds.", modelSizePercent, modelPickPercent)
	//}

	logging.Infof("Building cluster set from %d with a base seed count of %d (top %d of %d each iteration).", len(addrs), modelSize, pickCount, pickSize)
	logging.Infof("Density threshold is %f.", density)

	model := modeling.CreateClusteringModel(addrs)

	logging.Infof("Done generating model. Now writing to output path at '%s'.", outputPath)

	err = model.Save(outputPath)

	if err != nil {
		logging.ErrorF(err)
	}

	logging.Successf("Successfully generated model and wrote results to file '%s'.", outputPath)

}