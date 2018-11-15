package statemachine

import (
	"github.com/lavalamp-/ipv666/common/config"
	"log"
	"github.com/lavalamp-/ipv666/common/fs"
	"path/filepath"
	"fmt"
	"os"
	"github.com/lavalamp-/ipv666/common/data"
	"github.com/rcrowley/go-metrics"
)

var s3PushFileCounter = metrics.NewCounter()
var s3PushErrorCounter = metrics.NewCounter()
var s3BytesPushedCounter = metrics.NewCounter()

func init() {
	metrics.Register("s3push.s3_push.count", s3PushFileCounter)
	metrics.Register("s3push.s3_error.count", s3PushErrorCounter)
	metrics.Register("s3push.s3_push_bytes.count", s3BytesPushedCounter)
}

//TODO remove S3 functionality?

func pushFilesToS3(conf *config.Configuration) (error) {
	allDirs := conf.GetAllExportDirectories()
	log.Printf("Now starting to push all non-most-recent files from %d directories.", len(allDirs))
	for _, curDir := range allDirs {
		err := syncDirectoryWithS3(curDir, conf)
		if err != nil {
			log.Printf("Error thrown when syncing directory '%s' with S3: %e", curDir, err)
			if conf.ExitOnFailedSync {
				return err
			}
		}
	}
	log.Printf("All %d directories successfully exported to S3.", len(allDirs))
	return nil
}

func syncDirectoryWithS3(directory string, conf *config.Configuration) (error) {
	log.Printf("Processing content of directory '%s'.", directory)
	exportFiles, err := fs.GetNonMostRecentFilesFromDirectory(directory)
	if err != nil {
		s3PushErrorCounter.Inc(1)
		log.Printf("Error thrown when attempting to gather files for export in directory '%s'.", directory)
		return err
	} else if len(exportFiles) == 0 {
		log.Printf("No files found for export in directory '%s'.", directory)
		return nil
	}
	log.Printf("A total of %d files were found for export in directory '%s'.", len(exportFiles), directory)
	for _, curFileName := range exportFiles {
		curFilePath := filepath.Join(directory, curFileName)
		log.Printf("Now exporting file at local file path '%s' to THE GREAT BEYONDDDD.", curFilePath)
		err := syncFileWithS3(curFilePath, curFilePath, conf)
		if err != nil {
			log.Printf("Error thrown when attempting to sync file at path '%s' to S3: %e", curFilePath, err)
			if conf.ExitOnFailedSync {
				return err
			}
		}
	}
	log.Printf("All files in directory at '%s' processed.", directory)
	return nil
}

func syncFileWithS3(localPath string, remotePath string, conf *config.Configuration) (error) {
	zipFilePath := fmt.Sprintf("%s.zip", localPath)
	log.Printf("Zipping up file at path '%s' to file at path '%s'.", localPath, zipFilePath)
	err := fs.ZipFiles([]string{localPath}, zipFilePath)
	if err != nil {
		s3PushErrorCounter.Inc(1)
		log.Printf("Failed to zip up file at path '%s'. Stopping export.", localPath)
		zipErr := os.Remove(zipFilePath)
		if zipErr != nil {
			s3PushErrorCounter.Inc(1)
			log.Printf("Another error was thrown when trying to delete zip file at path '%s': %e", zipFilePath, err)
		}
		return err
	}
	log.Printf("Successfully created zip file at path '%s'.", zipFilePath)
	fileSize, err := fs.CountFileSize(zipFilePath)
	if err != nil {
		s3PushErrorCounter.Inc(1)
		log.Printf("Error thrown when getting file size of file at path '%s': %e", zipFilePath, err)
		zipErr := os.Remove(zipFilePath)
		if zipErr != nil {
			s3PushErrorCounter.Inc(1)
			log.Printf("Another error was thrown when trying to delete zip file at path '%s': %e", zipFilePath, err)
		}
		return err
	}
	log.Printf("Moving file at '%s' to S3 bucket (%s). %d total bytes.", zipFilePath, remotePath, fileSize)
	err = data.PushFileToS3FromConfig(zipFilePath, remotePath, conf)
	if err != nil {
		s3PushErrorCounter.Inc(1)
		log.Printf("Failed to move file at path '%s' to S3. Stopping export.", zipFilePath)
		zipErr := os.Remove(zipFilePath)
		if zipErr != nil {
			s3PushErrorCounter.Inc(1)
			log.Printf("Another error was thrown when trying to delete zip file at path '%s': %e", zipFilePath, err)
		}
		return err
	}
	log.Printf("Deleting zip file at '%s'.", zipFilePath)
	err = os.Remove(zipFilePath)
	if err != nil {
		s3PushErrorCounter.Inc(1)
		log.Printf("Error thrown when attempting to delete zip file at path '%s': %e", zipFilePath, err)
		return err
	}
	log.Printf("Successfully moved file at '%s' to S3 with compression.", localPath)
	s3PushFileCounter.Inc(1)
	s3BytesPushedCounter.Inc(fileSize)
	return nil
}
