package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"wavefy-be/config"
)

// NewR2Client creates an S3-compatible client for Cloudflare R2.
func NewR2Client(ctx context.Context, cfg config.R2Config) (*s3.Client, error) {
	awsCfg, err := awsConfig.LoadDefaultConfig(
		ctx,
		awsConfig.WithRegion(cfg.Region),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		),
		awsConfig.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
				if cfg.Endpoint == "" {
					return aws.Endpoint{}, &aws.EndpointNotFoundError{}
				}
				return aws.Endpoint{
					URL:               cfg.Endpoint,
					SigningRegion:     cfg.Region,
					HostnameImmutable: true,
				}, nil
			}),
		),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	}), nil
}
