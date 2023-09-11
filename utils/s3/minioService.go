package s3

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"mime/multipart"
	"os"
)

var MinioClient *minio.Client

func ConnectToMinio()  {
	var err error

	endpoint := os.Getenv("minioUrl")
	accessKey := os.Getenv("minioAccessKey")
	secretKey := os.Getenv("minioSecretKey")

	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	imageBucketName := os.Getenv("imageBucketName")
	imageBucketExists, err := MinioClient.BucketExists(context.Background(), imageBucketName)
	if err != nil {
		log.Fatal(err)
	}
	if !imageBucketExists {
		err = MinioClient.MakeBucket(context.Background(), imageBucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatal(err)
		}
	}
	mediaGroupBucketName := os.Getenv("mediaGroupBucketName")
	mediaGroupBucketExists ,err := MinioClient.BucketExists(context.Background(), mediaGroupBucketName)
	if err != nil {
		log.Fatal(err)
	}
	if !mediaGroupBucketExists {
		err = MinioClient.MakeBucket(context.Background(), mediaGroupBucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatal(err)
		}
	}
	videoBucketName := os.Getenv("videoBucketName")
	videoBucketExists ,err := MinioClient.BucketExists(context.Background(), videoBucketName)
	if err != nil {
		log.Fatal(err)
	}
	if !videoBucketExists {
		err = MinioClient.MakeBucket(context.Background(), videoBucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatal(err)
		}
	}
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

func GetFile(bucketName string, objectName string) (io.Reader, error) {
	var buffer bytes.Buffer
	reader ,err := MinioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	defer reader.Close()
	if err != nil {
		return nil, err
	}
	_,err = io.Copy(&buffer, reader)
	if err != nil {
		return nil, err
	}
	fileReader := &buffer
	return fileReader, nil
}

func DeleteFile(bucketName, objectName string) error {

	err := MinioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil

}

//func GetMultipartFile(ctx context.Context, bucketName, objectName string) (*multipart.File, error) {
//	reader, err := MinioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
//	if err != nil {
//		return
//	}
//	multipart.File(reade)
//}
