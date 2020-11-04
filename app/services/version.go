package services

import (
	"paltronus-backend/app/models"
)

func QueryAllVersions() (*[]models.Version, error) {
	versions := []models.Version{}
	if result := DB.Find(&versions); result.Error != nil {
		return &versions, result.Error
	}
	return &versions, nil
}

func QueryVersion(id int) (*models.Version, error) {
	version := models.Version{}
	if result := DB.First(&version, id); result.Error != nil {
		return &version, result.Error
	}
	return &version, nil
}

func InsertVersion(version models.Version) (*models.Version, error) {
	if result := DB.Create(&version); result.Error != nil {
		return &version, result.Error
	}
	return &version, nil
}

func QueryFileVersions(id int) (*[]models.Version, error) {
	versions := []models.Version{}
	if result := DB.Find(&versions, "file.id = ?", id); result.Error != nil {
		return &versions, result.Error
	}
	return &versions, nil
}