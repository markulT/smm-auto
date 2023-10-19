package s3

import (
	"context"
	"io"
	"mime/multipart"
	"os"
)

func LoadAudio(objectName string, file *multipart.File) error {
	audioBucketName := os.Getenv("audioBucketName")
	err := LoadMultipartFile(context.Background(), audioBucketName, objectName, file)
	if err != nil {
		return err
	}
	return nil
}

func GetAudio(objectName string) (io.Reader, error) {
	audioBucketName := os.Getenv("audioBucketName")
	reader , err := GetFile(audioBucketName, objectName)
	if err != nil {
		return nil,err
	}
	return reader, nil
}

func DeleteAudio(objectName string) error {
	audioBucketName := os.Getenv("audioBucketName")
	err := DeleteFile(audioBucketName, objectName)
	if err != nil {
		return err
	}
	return nil
}
