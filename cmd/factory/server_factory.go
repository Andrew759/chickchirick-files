package factory

import (
	"chickchirick-files/internal/controller/c_controller"
	"chickchirick-files/internal/controller/service/file"
	"chickchirick-files/pkg/chirick_config"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func BuildAndServe(s3Client *s3.Client) {
	err := BuildServer(s3Client)
	if err != nil {
		panic(err)
	}
}

func BuildServer(awsClient *s3.Client) error {
	e := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		viper.GetString(chirick_config.UserAPPUrl),
		viper.GetString(chirick_config.MessagesAppUrl),
	}
	e.Use(cors.New(config))

	InitFileServer(e, awsClient)

	err := e.Run()
	if err != nil {
		return err
	}

	return nil
}

func InitFileServer(e *gin.Engine, awsClient *s3.Client) {
	fileService := file.FileController{
		Controller: c_controller.Controller{
			E: e,
			DI: &c_controller.DIContainer{
				Client: *awsClient,
			},
		},
	}

	fileService.RegisterRoutes()
}
