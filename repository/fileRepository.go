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
	FindManyByIDList(c context.Context, fIDs []uuid.UUID) ([]models.File, error)
	DeleteManyByIdList(c context.Context, fIDs []uuid.UUID) error
}

type fileRepoImpl struct {}

func NewFileRepo() FileRepository {
	return &fileRepoImpl{}
}

func (fr *fileRepoImpl) DeleteManyByIdList(c context.Context, fIDs []uuid.UUID) error {
	fileCollection := utils.DB.Collection("files")
	_, err:= fileCollection.DeleteMany(c, bson.M{"_id":bson.M{"$in":fIDs}})
	if err != nil {
		return err
	}
	return nil
}

func (fr *fileRepoImpl) FindManyByIDList(c context.Context, fIDs []uuid.UUID) ([]models.File, error) {
	var fileList []models.File
	fileCollection := utils.DB.Collection("files")
	fc, err := fileCollection.Find(c, bson.M{"_id":bson.M{"$in":fIDs}})
	if err != nil {
		return nil, err
	}
	defer fc.Close(c)
	for fc.Next(c) {
		var file models.File
		err := fc.Decode(&file)
		if err != nil {
			return nil, err
		}
		fileList = append(fileList, file)
	}
	return fileList, nil
}

func (fr *fileRepoImpl) FindByID(c context.Context, fID uuid.UUID) (*models.File, error) {
	f := new(models.File)

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