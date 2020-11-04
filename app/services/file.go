package services

import (
	"paltronus-backend/app/models"
)

func QueryAllFiles() (*[]models.File, error) {
	files := []models.File{}
	if result := DB.Find(&files); result.Error != nil {
		return &files, result.Error
	}
	return &files, nil
}

func QueryFile(id int) (*models.File, error) {
	file := models.File{}
	if result := DB.First(&file, id); result.Error != nil {
		return &file, result.Error
	}
	return &file, nil
}

func InsertFile(file models.File) (*models.File, error) {
	if result := DB.Create(&file); result.Error != nil {
		return &file, result.Error
	}
	return &file, nil
}
