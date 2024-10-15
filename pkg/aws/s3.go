package aws

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/swibly/swibly-api/config"
)

func (svc *AWSService) UploadFileHeader(key string, file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("unable to open file: %v", err)
	}
	defer f.Close()

	_, err = svc.S3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(config.S3.Bucket),
		Key:         aws.String(key),
		Body:        f,
		ContentType: aws.String(file.Header.Get("Content-Type")),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("unable to upload file to S3: %v", err)
	}

	fileURL := fmt.Sprintf("%s/%s", config.S3.URL, key)

	return fileURL, nil
}

func (svc *AWSService) UploadFile(key string, file *os.File) (string, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("unable to get file info: %v", err)
	}

	_, err = svc.S3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(config.S3.Bucket),
		Key:           aws.String(key),
		Body:          file,
		ContentType:   aws.String("application/octet-stream"),
		ACL:           types.ObjectCannedACLPublicRead,
		ContentLength: aws.Int64(fileInfo.Size()),
	})
	if err != nil {
		return "", fmt.Errorf("unable to upload file to S3: %v", err)
	}

	fileURL := fmt.Sprintf("https://%s.%s/%s", config.S3.Bucket, config.S3.SURL, key)

	return fileURL, nil
}

func UploadProjectImage(file *multipart.FileHeader) {
}
