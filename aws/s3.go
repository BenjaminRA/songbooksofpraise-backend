package aws

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var awsS3Client *s3.Client

func GetS3Client() *s3.Client {
	if awsS3Client == nil {
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_S3_REGION")))
		if err != nil {
			panic(err)
		}

		awsS3Client = s3.NewFromConfig(cfg)
	}

	return awsS3Client
}

func S3UploadFile(path string, key string, bucket string) string {
	client := GetS3Client()

	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("File not found: ", path)

		panic(err)
	}

	uploader := manager.NewUploader(client)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(file),
	})

	if err != nil {
		fmt.Printf("Couldn't upload file %v to %v. Here's why: %v\n",
			path, bucket, err)

		panic(err)
	}

	return result.Location
}

func S3DeleteFile(key string, bucket string) error {
	client := GetS3Client()

	_, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return err
}
