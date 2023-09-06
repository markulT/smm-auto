package s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"
)

var MinioClient *minio.Client

func ConnectToMinio()  {
	fmt.Println("Connecting to minio")
	var err error

	endpoint := os.Getenv("minioUrl")
	accessKey := os.Getenv("minioAccessKey")
	secretKey := os.Getenv("minioSecretKey")

	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func CreateNewBucket(bucketName string)  {
	err := MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		log.Fatal(err)
		return
	}
}

func LoadFile(ctx context.Context, bucketName string, objectName string, file *os.File) error {
	fileInfo, err := file.Stat()
	buffer := make([]byte, fileInfo.Size())
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}
	_, err = MinioClient.PutObject(ctx, bucketName, objectName, bytes.NewReader(buffer), fileInfo.Size(), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return err
	}
	return nil
}

func LoadMultipartFile(ctx context.Context, bucketName string, objectName string, file *multipart.File) error {
	_, err := MinioClient.PutObject(context.Background(), bucketName, objectName, *file, -1, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func GetFile(ctx context.Context, bucketName, objectName string) (*os.File, error) {

	reader, err := MinioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var buffer strings.Builder
	_, err = io.Copy(&buffer, reader)
	if err != nil {
		return nil, err
	}
	file, err := os.CreateTemp("", "minio_file_")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = file.WriteString(buffer.String())
	if err != nil {
		return nil, err
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return file, nil
}
func GetMultipartFile(ctx context.Context, bucketName, objectName string) (*multipart.File, error) {
	reader, err := MinioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return
	}
	multipart.File(reade)
}
