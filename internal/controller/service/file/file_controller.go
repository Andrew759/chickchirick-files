package file

import (
	"chickchirick-files/internal/controller/c_controller"
	"chickchirick-files/pkg/chirick_config"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type FileController struct {
	c_controller.Controller
}

func (fc *FileController) RegisterRoutes() {
	e := fc.E

	e.POST("/file/:uuid", fc.UploadFile)
	e.GET("/file/:uuid", fc.GetFile)
	e.DELETE("/file/:uuid", fc.DeleteFile)
}

func (fc *FileController) UploadFile(c *gin.Context) {
	fileUuid := c.Param("uuid")

	file, fHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not found"})
		return
	}

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}(file)

	_, err = fc.DI.Client.PutObject(c.Request.Context(), &s3.PutObjectInput{
		Bucket: aws.String(viper.GetString(chirick_config.SeaweedMainBucket)),
		Key:    aws.String(fileUuid),
		Body:   file,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payload": "file upload success", "filename": fHeader.Filename})
}

func (fc *FileController) GetFile(c *gin.Context) {
	fileUuid := c.Param("uuid")

	output, err := fc.DI.Client.GetObject(c.Request.Context(), &s3.GetObjectInput{
		Bucket: aws.String(viper.GetString(chirick_config.SeaweedMainBucket)),
		Key:    aws.String(fileUuid),
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	defer output.Body.Close()

	contentType := "image/jpeg" // значение по умолчанию
	if output.ContentType != nil {
		contentType = *output.ContentType
	}

	contentLength := int64(0)
	if output.ContentLength != nil {
		contentLength = *output.ContentLength
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, output.Body, nil)
}

func (fc *FileController) DeleteFile(c *gin.Context) {
	fileUuid := c.Param("uuid")

	_, err := fc.DI.Client.DeleteObject(c.Request.Context(), &s3.DeleteObjectInput{
		Bucket: aws.String(viper.GetString(chirick_config.SeaweedMainBucket)),
		Key:    aws.String(fileUuid),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete file: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
