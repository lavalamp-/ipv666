package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"log"
	"github.com/lavalamp-/ipv666/common/fs"
	"path/filepath"
	"os"
	"github.com/rcrowley/go-metrics"
)


var cleanUpFileCounter = metrics.NewCounter()

func init() {
	metrics.Register("cleanup_file_count", cleanUpFileCounter)
}

func cleanUpNonRecentFiles(conf *config.Configuration) (error) {
	allDirs := conf.GetAllExportDirectories()
	log.Printf("Now starting to delete all non-recent files from %d directories.", len(allDirs))
	for _, curDir := range allDirs {
		log.Printf("Processing content of directory '%s'.", curDir)
		exportFiles, err := fs.GetNonMostRecentFilesFromDirectory(curDir)
		if err != nil {
			log.Printf("Error thrown when attempting to gather files for deletion in directory '%s'.", curDir)
			return err
		} else if len(exportFiles) == 0 {
			log.Printf("No files found for export in directory '%s'.", curDir)
			continue
		}
		for _, curFileName := range exportFiles {
			curFilePath := filepath.Join(curDir, curFileName)
			log.Printf("Deleting file at path '%s'.", curFilePath)
			err := os.Remove(curFilePath)
			if err != nil {
				log.Printf("Error thrown when attempting to delete file at path '%s': %e", curFilePath, err)
				return err
			}
			cleanUpFileCounter.Inc(1)
			log.Printf("Successfully deleted file at path '%s'.", curFilePath)
		}
		log.Printf("Deleted all files in directory '%s'.", curDir)
	}
	log.Printf("Successfully deleted all non-recent files from %d directories.", len(allDirs))
	return nil
}
