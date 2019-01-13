package statemachine

import (
	"github.com/lavalamp-/ipv666/internal/config"
	"github.com/lavalamp-/ipv666/internal/fs"
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/rcrowley/go-metrics"
	"os"
	"path/filepath"
)


var cleanUpFileCounter = metrics.NewCounter()

func init() {
	metrics.Register("cleanup.files.count", cleanUpFileCounter)
}

func cleanUpNonRecentFiles() error {
	allDirs := config.GetAllExportDirectories()
	logging.Infof("Now starting to delete all non-recent files from %d directories.", len(allDirs))
	for _, curDir := range allDirs {
		logging.Debugf("Processing content of directory '%s'.", curDir)
		exportFiles, err := fs.GetNonMostRecentFilesFromDirectory(curDir)
		if err != nil {
			logging.Warnf("Error thrown when attempting to gather files for deletion in directory '%s'.", curDir)
			return err
		} else if len(exportFiles) == 0 {
			logging.Debugf("No files found for export in directory '%s'.", curDir)
			continue
		}
		for _, curFileName := range exportFiles {
			curFilePath := filepath.Join(curDir, curFileName)
			logging.Debugf("Deleting file at path '%s'.", curFilePath)
			err := os.Remove(curFilePath)
			if err != nil {
				logging.Warnf("Error thrown when attempting to delete file at path '%s': %e", curFilePath, err)
				return err
			}
			cleanUpFileCounter.Inc(1)
			logging.Debugf("Successfully deleted file at path '%s'.", curFilePath)
		}
		logging.Debugf("Deleted all files in directory '%s'.", curDir)
	}
	logging.Infof("Successfully deleted all non-recent files from %d directories.", len(allDirs))
	return nil
}
