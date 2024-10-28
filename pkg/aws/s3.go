package aws

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
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

func (svc *AWSService) UploadFile(key string, file multipart.File) (string, error) {
	newKey := fmt.Sprintf("%s/%s-%d.webp", config.Router.Environment, key, time.Now().Unix())

	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		return "", fmt.Errorf("unable to read file: %v", err)
	}

	imgData := bytes.NewReader(buf.Bytes())
	img, _, err := image.Decode(imgData)
	if err != nil {
		return "", fmt.Errorf("unable to decode image: %v", err)
	}

	imgData.Seek(0, io.SeekStart)

	img = adjustOrientation(imgData, img)

	processedImgBuf := new(bytes.Buffer)
	err = webp.Encode(processedImgBuf, img, &webp.Options{
		Lossless: true,
	})
	if err != nil {
		return "", fmt.Errorf("unable to encode image to WebP: %v", err)
	}

	_, err = svc.s3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(config.S3.Bucket),
		Key:         aws.String(newKey),
		Body:        bytes.NewReader(processedImgBuf.Bytes()),
		ContentType: aws.String("image/webp"),
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

	src, err := file.Open()
	if err != nil {
		return "", ErrUnableToOpenFile
	}
	defer src.Close()

	outputPath := fmt.Sprintf("projects/%d", projectID)

	url, err := AWS.UploadFile(outputPath, src)
	if err != nil {
		return "", ErrUnableToUploadFile
	}

	return url, nil
}

func UploadUserImage(userID uint, file *multipart.FileHeader) (string, error) {
	const maxFileSize = 5 * 1024 * 1024

	if file.Size > maxFileSize {
		return "", ErrFileTooLarge
	}

	src, err := file.Open()
	if err != nil {
		return "", ErrUnableToOpenFile
	}
	defer src.Close()

	outputPath := fmt.Sprintf("users/%d", userID)

	url, err := AWS.UploadFile(outputPath, src)
	if err != nil {
		return "", ErrUnableToUploadFile
	}

	return url, nil
}

func DeleteProjectImage(filename string) error {
	return AWS.DeleteFile(fmt.Sprintf("projects/%s", filepath.Base(filename)))
}

func DeleteUserImage(filename string) error {
	return AWS.DeleteFile(fmt.Sprintf("users/%s", filepath.Base(filename)))
}

func adjustOrientation(reader io.Reader, img image.Image) image.Image {
	exifData, err := exif.Decode(reader)
	if err == nil {
		orientTag, err := exifData.Get(exif.Orientation)
		if err == nil {
			orient, err := orientTag.Int(0)
			if err == nil {
				switch orient {
				case 2:
					img = imaging.FlipH(img)
				case 3:
					img = imaging.Rotate180(img)
				case 4:
					img = imaging.FlipV(img)
				case 5:
					img = imaging.Transpose(img)
				case 6:
					img = imaging.Rotate270(img)
				case 7:
					img = imaging.Transverse(img)
				case 8:
					img = imaging.Rotate90(img)
				}
			}
		}
	}
	return img
}
