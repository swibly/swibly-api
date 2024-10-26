package aws

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/swibly/swibly-api/config"
)

var (
	ErrUnableToOpenFile    = fmt.Errorf("unable to open file")
	ErrUnableToDecode      = fmt.Errorf("unable to decode file")
	ErrUnableToEncode      = fmt.Errorf("unable to encode file")
	ErrUnsupportedFileType = fmt.Errorf("unsupported file type")
	ErrUnableToUploadFile  = fmt.Errorf("unable to upload file")
	ErrFileTooLarge        = fmt.Errorf("file too large")
)

func (svc *AWSService) UploadFile(key string, file io.Reader) (string, error) {
	newKey := fmt.Sprintf("%s/%s", config.Router.Environment, key)

	_, err := svc.s3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(config.S3.Bucket),
		Key:         aws.String(newKey),
		Body:        file,
		ContentType: aws.String("application/octet-stream"),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("unable to upload file to S3: %v", err)
	}

	fileURL := fmt.Sprintf("https://%s.%s/%s", config.S3.Bucket, config.S3.SURL, newKey)

	return fileURL, nil
}

func (svc *AWSService) DeleteFile(key string) error {
	newKey := fmt.Sprintf("%s/%s", config.Router.Environment, key)

	_, err := svc.s3.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(config.S3.Bucket),
		Key:    aws.String(newKey),
	})
	if err != nil {
		return fmt.Errorf("unable to delete file from S3: %v", err)
	}

	return nil
}

func UploadProjectImage(projectID uint, file *multipart.FileHeader) (string, error) {
	const maxFileSize = 5 * 1024 * 1024

	if file.Size > maxFileSize {
		return "", ErrFileTooLarge
	}

	ext := strings.ToLower(path.Ext(file.Filename))

	if !slices.Contains([]string{".png", ".jpg", ".jpeg"}, ext) {
		return "", ErrUnsupportedFileType
	}

	src, err := file.Open()
	if err != nil {
		return "", ErrUnableToOpenFile
	}
	defer src.Close()

	img, err := imaging.Decode(src)
	if err != nil {
		return "", ErrUnableToDecode
	}

	outputPath := fmt.Sprintf("projects/%d-%d.webp", time.Now().Unix(), projectID)
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

func DeleteProjectImage(filename string) error {
	return AWS.DeleteFile(fmt.Sprintf("projects/%s", filepath.Base(filename)))
}

func UploadUserImage(userID uint, file *multipart.FileHeader) (string, error) {
	const maxFileSize = 5 * 1024 * 1024

	if file.Size > maxFileSize {
		return "", ErrFileTooLarge
	}

	ext := strings.ToLower(path.Ext(file.Filename))

	if !slices.Contains([]string{".png", ".jpg", ".jpeg"}, ext) {
		return "", ErrUnsupportedFileType
	}

	src, err := file.Open()
	if err != nil {
		return "", ErrUnableToOpenFile
	}
	defer src.Close()

	img, err := imaging.Decode(src)
	if err != nil {
		return "", ErrUnableToDecode
	}

	outputPath := fmt.Sprintf("users/%d-%d.webp", time.Now().Unix(), userID)
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

func DeleteUserImage(filename string) error {
	return AWS.DeleteFile(fmt.Sprintf("users/%s", filepath.Base(filename)))
}
