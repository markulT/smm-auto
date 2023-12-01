package s3

import (
	"context"
	"github.com/google/uuid"
	"golearn/models"
	"io"
)

type FileService interface {
	GetFileByID(c context.Context, fID uuid.UUID) (io.Reader, error)
}

type FileRepository interface {
	Save(c context.Context,file *models.File) error
	DeleteByID(c context.Context, fID uuid.UUID) error
	FindByID(c context.Context, fID uuid.UUID) (*models.File, error)
	FindManyByIDList(c context.Context, fIDs []uuid.UUID) ([]models.File, error)
}

type fileServiceImpl struct {
	fileRepo FileRepository
}

func NewFileService(fileRepo FileRepository) FileService {
	fs := fileServiceImpl{fileRepo: fileRepo}
	return &fs
}

func (fs *fileServiceImpl) GetFileByID( c context.Context,fID uuid.UUID) (io.Reader, error) {
	fileModel, err := fs.fileRepo.FindByID(c, fID)
	if err != nil {
		return nil, err
	}
	return GetFile(fileModel.BucketName,fileModel.ID.String())
}


