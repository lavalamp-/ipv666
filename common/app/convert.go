package app

import (
	"github.com/lavalamp-/ipv666/common/addressing"
	"github.com/lavalamp-/ipv666/common/logging"
)

func RunConvert(inputPath string, outputPath string, outputType string) {

	logging.Infof("Reading IPv6 addresses from path at '%s', converting to type '%s', and writing results to '%s'.", inputPath, outputType, outputPath)

	var err error
	addrs, err := addressing.ReadIPsFromFile(inputPath)

	if err != nil {
		logging.ErrorF(err)
	}
	logging.Debugf("Successfully read %d addresses from file '%s'.", len(addrs), inputPath)

	switch outputType {
	case "txt":
		err = addressing.WriteIPsToHexFile(outputPath, addrs)
	case "bin":
		err = addressing.WriteIPsToBinaryFile(outputPath, addrs)
	case "hex":
		err = addressing.WriteIPsToFatHexFile(outputPath, addrs)
	}

	if err != nil {
		logging.ErrorF(err)
	}

	logging.Successf("Successfully wrote IP addresses to '%s' with file type of '%s'.", outputPath, outputType)

}
