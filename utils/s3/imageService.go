package s3

import (
	"context"
	"mime/multipart"
)

func LoadImage(ctx context.Context, objectName string, file *multipart.File) error {
	err := LoadMultipartFile(ctx, "image-bucket", objectName, file)
	if err != nil {
		return err
	}
	return nil
}
func DeleteImage()  {

}
func GetImage(ctx context.Context,objectName string) (multipart.File, error) {
	return nil, nil
}