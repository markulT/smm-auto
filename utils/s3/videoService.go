package s3

import (
	"context"
	"io"
	"mime/multipart"
	"os"
)

func LoadVideo(objectName string, file *os.File) error  {
	videoBucketName := os.Getenv("videoBucketName")
	err := LoadFile(context.Background(), videoBucketName, objectName, file )
	if err != nil {
		return err
	}
	return nil
}
func LoadVideoMultipart(objectName string, file *multipart.File) error {
	videoBucketName := os.Getenv("videoBucketName")
	err := LoadMultipartFile(context.Background(), videoBucketName, objectName, file)
	if err != nil {
		return err
	}
	return nil
}
func GetVideo(objectName string) (io.Reader, error) {
	videoBucketName := os.Getenv("videoBucketName")
	reader , err := GetFile(videoBucketName, objectName)
	if err != nil {
		return nil,err
	}
	return reader, nil
}
func DeleteVideo(objectName string) error {
	videoBucketName := os.Getenv("videoBucketName")
	err := DeleteFile(videoBucketName, objectName)
	if err != nil {
		return err
	}
	return nil
}
