package s3

import (
	"context"
	"io"
	"mime/multipart"
	"os"
)



func LoadMedia(ctx context.Context, objectName string, file *multipart.File) error {
	bucketName := os.Getenv("mediaGroupBucketName")
	err := LoadMultipartFile(ctx, bucketName, objectName, file)
	if err != nil {
		return err
	}
	return nil
}

func DeleteMedia(objectName string) error {
	bucketName := os.Getenv("mediaGroupBucketName")
	err := DeleteFile(bucketName, objectName)
	if err != nil {
		return err
	}
	return nil
}

func GetMedia(objectName string) (io.Reader, error) {
	bucketName := os.Getenv("mediaGroupBucketName")
	reader, err := GetFile(bucketName, objectName)
	if err != nil {
		return nil, err
	}
	return reader, nil
}