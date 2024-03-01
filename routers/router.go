package routers

import (
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/apis"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/registerModel", apis.RegisterModel)
	r.GET("/getModelInfo", apis.GetModelInfo)
	r.GET("/getModelInfo/:modelName", apis.GetModelInfoByName)
	r.POST("/uploadModel/:modelName", apis.UploadModel)
	r.GET("/downloadModel/:modelName/model.zip", apis.DownloadModel)
	return r
}
