package aws

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"path"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/chai2010/webp"
	"github.com/swibly/swibly-api/config"
)

var (
	ErrUnableToOpeFile     = fmt.Errorf("unable to open file")
	ErrUnableToDecode      = fmt.Errorf("unable to decode file")
	ErrUnableToEncode      = fmt.Errorf("unable to encode file")
	ErrUnsupportedFileType = fmt.Errorf("unsupported file type")
	ErrUnableToUploadFile  = fmt.Errorf("unable to upload file")
	ErrFileTooLarge        = fmt.Errorf("file too large")
)

func (svc *AWSService) UploadFile(key string, file io.Reader) (string, error) {
	_, err := svc.s3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(config.S3.Bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String("application/octet-stream"),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("unable to upload file to S3: %v", err)
	}

	fileURL := fmt.Sprintf("https://%s.%s/%s", config.S3.Bucket, config.S3.SURL, key)

	return fileURL, nil
}

func (svc *AWSService) DeleteFile(key string) error {
	_, err := svc.s3.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(config.S3.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("unable to delete file from S3: %v", err)
	}

	return nil
}

func UploadProjectImage(projectID uint, file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(path.Ext(file.Filename))

  if !slices.Contains([]string{".png", ".jpg", ".jpeg"}, ext) {
    return "", ErrUnsupportedFileType
  }

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("unable to open file: %v", err)
	}
	defer src.Close()

	var img image.Image
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(src)
		if err != nil {
			return "", ErrUnableToDecode
		}
	case ".png":
		img, err = png.Decode(src)
		if err != nil {
			return "", ErrUnableToDecode
		}
	default:
		return "", ErrUnsupportedFileType
	}

	outputPath := fmt.Sprintf("projects/%d.webp", projectID)
	var buf bytes.Buffer

	err = webp.Encode(&buf, img, nil)
	if err != nil {
		return "", ErrUnableToEncode
	}

	url, err := AWS.UploadFile(outputPath, &buf)
	if err != nil {
		return "", ErrUnableToUploadFile
	}

	return url, nil
}

func DeleteProjectImage(projectID uint) error {
	return AWS.DeleteFile(fmt.Sprintf("projects/%d.webp", projectID))
}
