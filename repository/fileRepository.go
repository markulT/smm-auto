package repository

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golearn/models"
	"golearn/utils"
)

type FileRepository interface {
	Save(c context.Context,file *models.File) error
	DeleteByID(c context.Context, fID uuid.UUID) error
	FindByID(c context.Context, fID uuid.UUID) (*models.File, error)
}

type fileRepoImpl struct {}

func NewFileRepo() FileRepository {
	return &fileRepoImpl{}
}

func (fr *fileRepoImpl) FindByID(c context.Context, fID uuid.UUID) (*models.File, error) {
	var f *models.File

	fileCollection := utils.DB.Collection("files")
	res := fileCollection.FindOne(c, bson.M{"_id":fID})

	if res.Err()!=nil {
		return nil, res.Err()
	}
	if err:=res.Decode(f);err!=nil {
		return nil, err
	}
	return f, nil
}

func (fr *fileRepoImpl) Save(c context.Context,f *models.File) error {
	fileCollection := utils.DB.Collection("files")
	_, err := fileCollection.InsertOne(c, f)
	if err != nil {
		return err
	}

	return nil
}

func (fr *fileRepoImpl) DeleteByID(c context.Context, fID uuid.UUID) error {

	fileCollection := utils.DB.Collection("files")
	_,err := fileCollection.DeleteOne(c, bson.M{"_id":fID})
	if err != nil {
		return err
	}
	return nil
}