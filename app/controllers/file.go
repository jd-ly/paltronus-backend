package controllers

import (
	"net/http"
	"time"

	"github.com/revel/revel"

	"paltronus-backend/app/models"
	"paltronus-backend/app/services"
)

type File struct {
	App
}

// GetFiles is used to get all file json
func (c *File) GetFiles() revel.Result {
	file, err := services.QueryAllFiles()
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(map[string]string{"status": "Invalid Request"})
	}
	return c.RenderJSON(file)
}

func (c *File) CreateFile() revel.Result {
	var jsonData map[string]interface{}
	c.Params.BindJSON(&jsonData)

	title, oktitle := jsonData["title"].(string)
	description, okdesc := jsonData["description"].(string)
	createdBy, okcreated := jsonData["createdBy"].(string)

	if !oktitle || !okdesc || !okcreated {
		return c.RenderJSON(map[string]string{"status": "Invalid Parameters"})
	}

	file := models.File{
		Title:        title,
		Description: description,
		CreatedBy:       createdBy,
		CreationDate:    time.Now().Format("2006-01-02 15:04:05"),
	}

	file.Validate(c.Validation)
	if c.Validation.HasErrors() {
		return c.RenderJSON(map[string]string{"status": "Invalid Parameters"})
	}

	_, err := services.InsertFile(file)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(map[string]string{"status": "Invalid Request"})
	}

	return c.RenderJSON(jsonData)
}

func (c *File) GetFile(id int) revel.Result {
	file, err := services.QueryFile(id)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(map[string]string{"status": "Invalid Request"})
	}
	return c.RenderJSON(file)
}

func (c *File) UpdateFile(id int) revel.Result {
	file, err := services.QueryFile(id)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(map[string]string{"status": "Invalid Request"})
	}
	return c.RenderJSON(file)
}

func (c *File) DeleteFile(id int) revel.Result {
	file, err := services.QueryFile(id)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(map[string]string{"status": "Invalid Request"})
	}
	return c.RenderJSON(file)
}