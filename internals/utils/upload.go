package utils

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

const S3Region = "us-east-1"
const S3Bucket = "novelism_bucket"

func UploadFileToS3(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(S3Region),
		Credentials: credentials.NewCredentials(&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     viper.GetString("aws.access_key"),
				SecretAccessKey: viper.GetString("aws.secret_access_key"),
			},
		}),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}

	svc := s3.New(sess)

	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	object := s3.PutObjectInput{
		Bucket:               aws.String(S3Bucket),
		Key:                  aws.String(fileHeader.Filename),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	}

	_, err = svc.PutObject(&object)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	return fileHeader.Filename, nil
}
