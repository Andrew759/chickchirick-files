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

	e.POST("/file:uuid", fc.FileUpload)
}

func (fc *FileController) FileUpload(c *gin.Context) {
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

	_, err = fc.DI.Client.PutObject(c, &s3.PutObjectInput{
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
