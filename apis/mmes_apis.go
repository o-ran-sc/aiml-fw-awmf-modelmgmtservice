/*
==================================================================================
Copyright (c) 2023 Samsung Electronics Co., Ltd. All Rights Reserved.

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
package apis

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/core"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/logging"
	"github.com/gin-gonic/gin"
)

type ModelInfo struct {
	ModelName string                 `json:"model-name"`
	RAppId    string                 `json:"rapp-id"`
	Metainfo  map[string]interface{} `json:"meta-info"`
}

type ModelInfoResponse struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MmeApiHandler struct {
	dbmgr core.DBMgr
}

func NewMmeApiHandler(dbMgr core.DBMgr) *MmeApiHandler {
	handler := &MmeApiHandler{
		dbmgr: dbMgr,
	}
	return handler
}

func (m *MmeApiHandler) RegisterModel(cont *gin.Context) {
	var returnCode int = http.StatusCreated
	var responseMsg string = "Model registered successfully"

	logging.INFO("Creating model...")
	bodyBytes, _ := io.ReadAll(cont.Request.Body)

	var modelInfo ModelInfo
	//Need to unmarshal JSON to Struct, to access request
	//data such as model name, rapp id etc
	err := json.Unmarshal(bodyBytes, &modelInfo)
	if err != nil || modelInfo.ModelName == "" {
		logging.ERROR("Error in unmarshalling")
		cont.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": string("Can not parse input data, provide mandatory details"),
		})
	} else {
		logging.INFO(modelInfo.ModelName, modelInfo.RAppId, modelInfo.Metainfo)
		modelInfoBytes, _ := json.Marshal(modelInfo)

		err := m.dbmgr.CreateBucket(modelInfo.ModelName)
		if err == nil {
			m.dbmgr.UploadFile(modelInfoBytes, modelInfo.ModelName+os.Getenv("INFO_FILE_POSTFIX"), modelInfo.ModelName)
		} else {
			returnCode = http.StatusInternalServerError
			responseMsg = err.Error()
		}
		cont.JSON(returnCode, gin.H{
			"code":    returnCode,
			"message": responseMsg,
		})
	}
}

/*
This API retrieves model info list managed in modelmgmtservice
*/
func (m *MmeApiHandler) GetModelInfo(cont *gin.Context) {
	logging.INFO("List all model API")
	bucketList, err := m.dbmgr.ListBucket(os.Getenv("INFO_FILE_POSTFIX"))
	if err != nil {
		statusCode := http.StatusInternalServerError
		logging.ERROR("Error occurred, send status code: ", statusCode)
		cont.JSON(statusCode, gin.H{
			"code":    statusCode,
			"message": "Unexpected Error in server, you can't get model information list",
		})
		return
	}

	modelInfoListResp := []ModelInfoResponse{}
	for _, bucket := range bucketList {
		modelInfoListResp = append(modelInfoListResp, ModelInfoResponse{
			Name: bucket.Name,
			Data: string(bucket.Object),
		})
	}

	cont.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": modelInfoListResp,
	})
}

/*
Provides the model details by param model name
*/
func (m *MmeApiHandler) GetModelInfoByName(cont *gin.Context) {
	logging.INFO("Get model info by name API ...")
	modelName := cont.Param("modelName")

	bucketObj := m.dbmgr.GetBucketObject(modelName+os.Getenv("INFO_FILE_POSTFIX"), modelName)
	modelInfoListResp := ModelInfoResponse{
		Name: modelName,
		Data: string(bucketObj),
	}

	cont.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": modelInfoListResp,
	})
}

// API to upload the trained model in zip format
// TODO : Model version as input

func (m *MmeApiHandler) UploadModel(cont *gin.Context) {
	logging.INFO("Uploading model API ...")
	modelName := cont.Param("modelName")
	//TODO convert multipart.FileHeader to []byted
	fileHeader, _ := cont.FormFile("file")
	//TODO : Accept only .zip file for trained model
	file, _ := fileHeader.Open()
	//TODO: Handle error response
	defer file.Close()
	byteFile, _ := io.ReadAll((file))

	logging.INFO("Uploading model : ", modelName)
	m.dbmgr.UploadFile(byteFile, modelName+os.Getenv("MODEL_FILE_POSTFIX"), modelName)
	cont.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": string("Model uploaded successfully.."),
	})
}

/*
API to download the trained model from  bucket
Input: model name in path params as "modelName"
*/
func (m *MmeApiHandler) DownloadModel(cont *gin.Context) {
	logging.INFO("Download model API ...")
	modelName := cont.Param("modelName")
	fileName := modelName + os.Getenv("MODEL_FILE_POSTFIX")
	fileByes := m.dbmgr.GetBucketObject(fileName, modelName)

	//Return file in api reponse using byte slice
	cont.Header("Content-Disposition", "attachment;"+fileName)
	cont.Header("Content-Type", "application/zip")
	cont.Data(http.StatusOK, "application/octet", fileByes)
}

func (m *MmeApiHandler) GetModel(cont *gin.Context) {
	logging.INFO("Fetching model")
	cont.IndentedJSON(http.StatusOK, " ")
}

func (m *MmeApiHandler) UpdateModel() {
	logging.INFO("Updating model...")
}

func (m *MmeApiHandler) DeleteModel() {
	logging.INFO("Deleting model...")
}
