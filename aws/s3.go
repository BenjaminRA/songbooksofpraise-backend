package aws

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

var awsS3Client *s3.Client

func GetS3Client() *s3.Client {
	if awsS3Client == nil {
		// Load environment variables from .env file
		_ = godotenv.Load()

		// Get AWS credentials from environment variables (using your custom names)
		accessKey := os.Getenv("AWS_S3_ACCESS_KEY")
		secretKey := os.Getenv("AWS_S3_SECRET_ACCESS_KEY")
		region := os.Getenv("AWS_S3_REGION")

		// Debug: Check if credentials are loaded
		if accessKey == "" || secretKey == "" || region == "" {
			panic(fmt.Errorf("AWS credentials not found. Please check your .env file:\n"+
				"AWS_S3_ACCESS_KEY: %s\n"+
				"AWS_S3_SECRET_ACCESS_KEY: %s\n"+
				"AWS_S3_REGION: %s",
				maskString(accessKey), maskString(secretKey), region))
		}

		// Force use of static credentials to avoid IMDS
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		)

		if err != nil {
			panic(fmt.Errorf("failed to load AWS config: %w", err))
		}

		awsS3Client = s3.NewFromConfig(cfg)
	}

	return awsS3Client
}

// Helper function to mask credentials for logging
func maskString(s string) string {
	if len(s) == 0 {
		return "[EMPTY]"
	}
	if len(s) <= 4 {
		return "[SET]"
	}
	return s[:4] + "****"
}

func S3UploadFile(path string, key string, bucket string) (string, error) {
	// Validate required environment variables
	if os.Getenv("AWS_S3_REGION") == "" {
		return "", fmt.Errorf("AWS_S3_REGION environment variable is required")
	}

	client := GetS3Client()

	file, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("file not found: %s - %w", path, err)
	}

	client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	uploader := manager.NewUploader(client)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(file),
	})

	if err != nil {
		return "", fmt.Errorf("couldn't upload file %s to %s: %w", path, bucket, err)
	}

	return result.Location, nil
}

func S3DeleteFile(key string, bucket string) error {
	client := GetS3Client()

	_, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return err
}
