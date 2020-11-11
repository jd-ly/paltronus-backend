package controllers

import (
	"net/http"
	"time"

	"github.com/revel/revel"

	"paltronus-backend/app/models"
	"paltronus-backend/app/services"
)

type Version struct {
	App
}

func (c *Version) GetVersions() revel.Result {
	version, err := services.QueryAllVersions()
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(map[string]string{"status": "Invalid Request"})
	}
	return c.RenderJSON(version)
}

func (c *Version) CreateVersion() revel.Result {
	var jsonData map[string]interface{}
	c.Params.BindJSON(&jsonData)

	rawData, okdata := jsonData["rawData"].(string)
	createdBy, okcreated := jsonData["createdBy"].(string)
	fileId, okfile := jsonData["file"].(int)

	if !okdata || ! okcreated || !okfile {
		return c.RenderJSON(map[string]string{"status": "Invalid Parameters"})
	}

	version := models.Version{
		RawData: 		rawData,
		CreatedBy:      createdBy,
		CreationDate:   time.Now().Format("2006-01-02 15:04:05"),
		File: 			fileId,
	}

	version.Validate(c.Validation)
	if c.Validation.HasErrors() {
		return c.RenderJSON(map[string]string{"status": "Invalid Parameters"})
	}

	_, erro := services.InsertVersion(version)
	if erro != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(map[string]string{"status": "Invalid Request"})
	}

	return c.RenderJSON(jsonData)
}

func (c *Version) GetVersion(id int) revel.Result {
	version, err := services.QueryVersion(id)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(map[string]string{"status": "Invalid Request"})
	}
	return c.RenderJSON(version)
}

func (c *Version) GetFileVersions(id int) revel.Result {
	version, err := services.QueryFileVersions(id)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(map[string]string{"status": "Invalid Request"})
	}
	return c.RenderJSON(version)
}

func (c *Version) GetLastVersion(id int) revel.Result {
	version, err := services.QueryLastVersion(id)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(map[string]string{"status": "Invalid Request"})
	}
	return c.RenderJSON(version)
}