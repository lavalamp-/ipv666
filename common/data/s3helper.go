package data

import (
	"github.com/lavalamp-/ipv666/common/config"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"log"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

func PushFileToS3(
	localPath,
	remotePath,
	bucketName,
	awsRegion,
	awsAccessKey,
	awsSecretKey string,
	) (error) {
	// TODO consider caching the client
	sess, err := session.NewSession(&aws.Config{
		Region:			aws.String(awsRegion),
		Credentials:	credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
	})
	if err != nil {
		log.Printf("Error thrown when establishing S3 session: %e", err)
		return err
	}
	uploader := s3manager.NewUploader(sess)
	localFile, err := os.Open(localPath)
	if err != nil {
		log.Printf("Error thrown when opening local file at path '%s' for S3 upload: %e", localPath, err)
		return err
	}
	defer localFile.Close()
	log.Printf("Now attempting to upload local file at path '%s' to bucket '%s' under key '%s'.", localPath, bucketName, remotePath)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(remotePath),
		Body: localFile,
	})
	if err != nil {
		log.Printf("Error thrown when uploading file at path '%s' to S3: %e", localPath, err)
		return err
	}
	log.Printf("Successfully uploaded local file at path '%s' to '%s'.", localPath, aws.StringValue(&result.Location))
	return nil
}

func PushFileToS3FromConfig(localPath string, remotePath string, conf *config.Configuration) (error) {
	return PushFileToS3(
		localPath,
		remotePath,
		conf.AWSBucketName,
		conf.AWSBucketRegion,
		conf.AWSAccessKey,
		conf.AWSSecretKey,
	)
}
