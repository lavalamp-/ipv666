package app

import (
	"github.com/lavalamp-/ipv666/internal/fs"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/lavalamp-/ipv666/internal/modeling"
)

func RunModelgen(inputPath string, outputPath string) {

	logging.Infof("Reading source addresses from file at path '%s'.", inputPath)

	var err error
	addrs, err := fs.ReadIPsFromFile(inputPath)

	if err != nil {
		logging.ErrorF(err)
	}
	logging.Debugf("Successfully read %d addresses from file '%s'.", len(addrs), inputPath)

	logging.Infof("Building cluster set from %d addresses.", len(addrs))

	model := modeling.CreateClusteringModel(addrs)

	logging.Infof("Done generating model. Now writing to output path at '%s'.", outputPath)

	err = model.Save(outputPath)

	if err != nil {
		logging.ErrorF(err)
	}

	logging.Successf("Successfully generated model and wrote results to file '%s'.", outputPath)

}