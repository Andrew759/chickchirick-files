package main

import (
	"chickchirick-files/cmd/factory"
	"chickchirick-files/pkg/chirick_config"
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

func main() {
	factory.InitViper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(viper.GetString(chirick_config.SeaweedRegion)),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			viper.GetString(chirick_config.SeaweedFSKey),
			viper.GetString(chirick_config.SeaweedFSSECRET),
			"",
		),
		),
	)
	if err != nil {
		//TODO: заменить на panic?
		log.Fatalf("failed to load s3 config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(viper.GetString(chirick_config.SeaweedFSUrl))
		o.UsePathStyle = true
	})

	_, _ = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(viper.GetString(chirick_config.SeaweedMainBucket)),
	})

	factory.BuildAndServe(s3Client)

}
