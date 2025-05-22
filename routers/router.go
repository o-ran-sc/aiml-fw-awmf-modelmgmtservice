/*
==================================================================================
Copyright (c) 2024 Samsung Electronics Co., Ltd. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
==================================================================================
*/
package routers

import (
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/apis"
	"github.com/gin-gonic/gin"
)

func InitRouter(handler *apis.MmeApiHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	api := r.Group("/ai-ml-model-registration/v1")
	{
		api.POST("/model-registrations", handler.RegisterModel)
		api.POST("/model-registrations/updateArtifact/:modelname/:modelversion/:artifactversion", handler.UpdateArtifact)
		api.GET("/model-registrations/:modelRegistrationId", handler.GetModelInfoById)
		api.PUT("/model-registrations/:modelRegistrationId", handler.UpdateModel)
		api.DELETE("/model-registrations/:modelRegistrationId", handler.DeleteModel)
		api.GET("/getModelInfo/:modelName", handler.GetModelInfoByName)
		api.POST("/uploadModel/:modelName/:modelVersion/:artifactVersion", handler.UploadModel)
		api.GET("/downloadModel/:modelName/:modelVersion/:artifactVersion/model.zip", handler.DownloadModel)
	}
	// As per R1-AP v6

	modelDiscovery := r.Group("/ai-ml-model-discovery/v1")
	{
		modelDiscovery.GET("/models", handler.GetModelInfo)
	}
	return r
}
