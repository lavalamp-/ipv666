package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"log"
	"github.com/lavalamp-/ipv666/common/fs"
	"path/filepath"
	"fmt"
	"os"
	"github.com/lavalamp-/ipv666/common/data"
)

func pushFilesToS3(conf *config.Configuration) (error) {
	// TODO break this down into multiple functions
	allDirs := conf.GetAllExportDirectories()
	log.Printf("Now starting to push all non-most-recent files from %d directories.", len(allDirs))
	for _, curDir := range allDirs {
		log.Printf("Processing content of directory '%s'.", curDir)
		exportFiles, err := fs.GetNonMostRecentFilesFromDirectory(curDir)
		if err != nil {
			log.Printf("Error thrown when attempting to gather files for export in directory '%s'.", curDir)
			return err
		} else if len(exportFiles) == 0 {
			log.Printf("No files found for export in directory '%s'.", curDir)
			continue
		}
		log.Printf("A total of %d files were found for export in directory '%s'.", len(exportFiles), curDir)
		for _, curFileName := range exportFiles {
			curFilePath := filepath.Join(curDir, curFileName)
			zipFilePath := fmt.Sprintf("%s.zip", curFilePath)
			log.Printf("Now exporting file at local file path '%s' to THE GREAT BEYONDDDD.", curFilePath)
			log.Printf("Zipping up file at path '%s' to file at path '%s'.", curFilePath, zipFilePath)
			err := fs.ZipFiles([]string{curFilePath}, zipFilePath)
			if err != nil {
				log.Printf("Failed to zip up file at path '%s'. Stopping export.", curFilePath)
				zipErr := os.Remove(zipFilePath)
				if zipErr != nil {
					log.Printf("Another error was thrown when trying to delete zip file at path '%s': %e", zipFilePath, err)
				}
				return err
			}
			log.Printf("Successfully created zip file at path '%s'.", zipFilePath)
			log.Printf("Moving file at '%s' to S3 bucket.", zipFilePath)
			err = data.PushFileToS3FromConfig(zipFilePath, zipFilePath, conf)
			if err != nil {
				log.Printf("Failed to move file at path '%s' to S3. Stopping export.", zipFilePath)
				zipErr := os.Remove(zipFilePath)
				if zipErr != nil {
					log.Printf("Another error was thrown when trying to delete zip file at path '%s': %e", zipFilePath, err)
				}
				return err
			}
			log.Printf("Deleting zip file at '%s'.", zipFilePath)
			err = os.Remove(zipFilePath)
			if err != nil {
				log.Printf("Error thrown when attempting to delete zip file at path '%s': %e", zipFilePath, err)
				return err
			}
			log.Printf("Successfully moved file at '%s' to S3 with compression.", curFilePath)
		}
		log.Printf("All files in directory at '%s' processed.", curDir)
	}
	log.Printf("All %d directories successfully exported to S3.", len(allDirs))
	return nil
}
