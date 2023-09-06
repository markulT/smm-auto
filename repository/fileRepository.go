package repository

import (
	"context"
	"golearn/models"
	"golearn/utils"
)

func SaveFile(file *models.File) error {
	//stmt, err := utils.DB.Preparex("insert into files (bucket_name, type, post_id) values ($1, $2, $3)")
	//if err != nil {
	//	return err
	//}
	//defer stmt.Close()
	//err = stmt.Get(file.BucketName, file.Type, file.PostID)
	//return err
	fileCollection := utils.DB.Collection("file")
	_, err := fileCollection.InsertOne(context.TODO(), file)
	if err != nil {
		return err
	}
	return nil
}
