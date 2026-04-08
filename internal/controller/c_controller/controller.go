package c_controller

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

type DIContainer struct {
	Client s3.Client
}

type Controller struct {
	E  *gin.Engine
	DI *DIContainer
}

type RequestHandler interface {
	RegisterRoutes()
}

type ControllerInterface interface {
	RequestHandler
}
