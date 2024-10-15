package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	c "github.com/swibly/swibly-api/config"
)

type AWSService struct {
	s3 *s3.Client
}

var AWS *AWSService

func NewAWSService() error {
	fmt.Println("Access Key:", c.S3.Access)
	fmt.Println("Secret Key:", c.S3.Secret)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(c.S3.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			c.S3.Access,
			c.S3.Secret,
			"",
		)),
	)
	if err != nil {
		return fmt.Errorf("unable to load AWS SDK config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.EndpointResolver = s3.EndpointResolverFromURL(c.S3.URL)
		o.UsePathStyle = true
	})

	AWS = &AWSService{s3: s3Client}

	return nil
}
