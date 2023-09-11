package s3

import (
	"context"
	"io"
	"mime/multipart"
	"os"
)

func LoadImage(ctx context.Context, objectName string, file *multipart.File) error {
	imageBucketName := os.Getenv("imageBucketName")
	err := LoadMultipartFile(ctx, imageBucketName, objectName, file)
	if err != nil {
		return err
	}
	return nil
}

func DeleteImage(objectName string) error {
	imageBucketName := os.Getenv("imageBucketName")
	err := DeleteFile(imageBucketName, objectName)
	if err != nil {
		return err
	}
	return nil
}

func GetImage(objectName string) (io.Reader, error) {
	imageBucketName := os.Getenv("imageBucketName")
	reader, err := GetFile(imageBucketName, objectName)
	if err != nil {
		return nil, err
	}
	return reader, nil
}