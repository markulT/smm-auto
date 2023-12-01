package s3

import (
	"context"
	"github.com/google/uuid"
	"golearn/models"
	"io"
)

type FileService interface {
	GetFileByID(c context.Context, fID uuid.UUID) (io.Reader, error)
	DeleteOneByID(c context.Context, fID uuid.UUID) error
	DeleteManyByID(c context.Context, fIDs []uuid.UUID) error
}

type FileRepository interface {
	Save(c context.Context,file *models.File) error
	DeleteByID(c context.Context, fID uuid.UUID) error
	FindByID(c context.Context, fID uuid.UUID) (*models.File, error)
	FindManyByIDList(c context.Context, fIDs []uuid.UUID) ([]models.File, error)
	DeleteManyByIdList(c context.Context, fIDs []uuid.UUID) error
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

func (fs *fileServiceImpl) DeleteOneByID(c context.Context, fID uuid.UUID) error {

	fileModel, err := fs.fileRepo.FindByID(c, fID)
	if err != nil {
		return err
	}
	_ = DeleteFile(fileModel.BucketName, fileModel.ID.String())
	//if err != nil {
	//	return err
	//}
	err = fs.fileRepo.DeleteByID(c, fID)
	if err != nil {
		return err
	}
	return nil
}

func (fs *fileServiceImpl) DeleteManyByID(c context.Context, fIDs []uuid.UUID) error {
	fileModels, err := fs.fileRepo.FindManyByIDList(c,fIDs)
	if err != nil {
		return err
	}
	for _, file := range fileModels {
		_ = DeleteFile(file.BucketName, file.ID.String())
	}
	err=fs.fileRepo.DeleteManyByIdList(c, fIDs)
	if err != nil {
		return err
	}
	return nil
}