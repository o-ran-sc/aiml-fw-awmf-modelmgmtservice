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
	"io"
	"net/http"
	"os"
	"fmt"

	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/core"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/db"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/logging"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

)

type MmeApiHandler struct {
	dbmgr core.DBMgr
	iDB   db.IDB
}

func NewMmeApiHandler(dbMgr core.DBMgr, iDB db.IDB) *MmeApiHandler {
	handler := &MmeApiHandler{
		dbmgr: dbMgr,
		iDB:   iDB,
	}
	return handler
}

func (m *MmeApiHandler) RegisterModel(cont *gin.Context) {

	var modelInfo models.ModelInfo

	if err := cont.ShouldBindJSON(&modelInfo); err != nil {
		cont.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id := uuid.New()
	modelInfo.Id = id.String()

	// TODO: validate the object

	if err := m.iDB.Create(modelInfo); err != nil {
		logging.ERROR("error", err)
		return
	}

	logging.INFO("model is saved.")

	cont.JSON(http.StatusCreated, gin.H{
		"modelInfo": modelInfo,
	})

	// logging.INFO("Creating model...")
	// bodyBytes, _ := io.ReadAll(cont.Request.Body)

	// var modelInfo models.ModelInfo
	// //Need to unmarshal JSON to Struct, to access request
	// //data such as model name, rapp id etc
	// err := json.Unmarshal(bodyBytes, &modelInfo)
	// if err != nil || modelInfo.ModelId.ModelName == "" {
	// 	logging.ERROR("Error in unmarshalling")
	// 	cont.JSON(http.StatusBadRequest, gin.H{
	// 		"code":    http.StatusBadRequest,
	// 		"message": string("Can not parse input data, provide mandatory details"),
	// 	})
	// } else {
	// 	id := uuid.New()
	// 	modelInfo.Id = id.String()
	// 	modelInfoBytes, _ := json.Marshal(modelInfo)
	// 	err := m.dbmgr.CreateBucket(modelInfo.ModelId.ModelName)
	// 	if err == nil {
	// 		m.dbmgr.UploadFile(modelInfoBytes, modelInfo.ModelId.ModelName+os.Getenv("INFO_FILE_POSTFIX"), modelInfo.ModelId.ModelName)
	// 	} else {
	// 		cont.JSON(http.StatusInternalServerError, gin.H{
	// 			"code":    http.StatusInternalServerError,
	// 			"message": err.Error(),
	// 		})
	// 	}
	// 	cont.JSON(http.StatusCreated, gin.H{
	// 		"modelinfo": modelInfoBytes,
	// 	})
	// }
}

/*
This API retrieves model info list managed in modelmgmtservice
*/
func (m *MmeApiHandler) GetModelInfo(cont *gin.Context) {

	logging.INFO("Get model info ")
	queryParams := cont.Request.URL.Query()
	//to check only modelName and modelVersion can be passed.
	allowedParams := map[string]bool{
		"modelName": true,
		"modelVersion": true,
	}

	for key := range queryParams {
		if !allowedParams[key] {
			cont.JSON(http.StatusBadRequest, gin.H{
				"error": "Only modelName and modelVersion are allowed",
			})
			return
		}
	}

	modelName:= cont.Query("modelName")
	modelVersion:= cont.Query("modelVersion")

	if modelName == "" {
		//return all modelinfo stored 

		models, err := m.iDB.GetAll()
		if err != nil {
			logging.ERROR("error:", err)
			cont.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error(),
			})
			return
		}
		cont.JSON(http.StatusOK, models)
		return
	} else {
		if modelVersion == "" {
			// get all modelInfo by model name
			modelInfos, err:= m.iDB.GetModelInfoByName(modelName)
			if err != nil {
				statusCode := http.StatusInternalServerError
				logging.ERROR("Error occurred, send status code: ", statusCode)
				cont.JSON(statusCode, gin.H{
					"code":    statusCode,
					"message": "Unexpected Error in server, you can't get model information list",
				})
				return
			}
			//to check record not found
			if len(modelInfos)==0{
				statusCode := http.StatusNotFound
				errMessage := fmt.Sprintf("Record not found with modelName: %s", modelName)
				logging.ERROR("Record not found, send status code: ", statusCode)
				cont.JSON(statusCode, gin.H{
					"code":    statusCode,
					"message": errMessage,
				})
				return
			}

			cont.JSON(http.StatusOK, gin.H{
				"modelinfoList":modelInfos,
			})
			return

		} else
		{
			// get all modelInfo by model name and version
			modelInfo, err:= m.iDB.GetModelInfoByNameAndVer(modelName, modelVersion)
			if err != nil {
				statusCode := http.StatusInternalServerError
				logging.ERROR("Error occurred, send status code: ", statusCode)
				cont.JSON(statusCode, gin.H{
					"code":    statusCode,
					"message": "Unexpected Error in server, you can't get model information list",
				})
				return
			}
			if modelInfo.Id == ""{
				statusCode := http.StatusNotFound
				errMessage := fmt.Sprintf("Record not found with modelName: %s and modelVersion: %s", modelName, modelVersion)
				logging.ERROR("Record not found, send status code: ", statusCode)
				cont.JSON(statusCode, gin.H{
					"code":    statusCode,
					"message": errMessage,
				})
				return
			}

			cont.JSON(http.StatusOK, gin.H{
				"modelinfo":modelInfo,
			})
			return
		}
	}
}

func (m *MmeApiHandler) GetModelInfoById(cont *gin.Context) {
	logging.INFO("Get model info by id ...")
	id := cont.Param("id")
	modelInfo, err := m.iDB.GetModelInfoById(id)
	if err != nil {
		logging.ERROR("error:", err)
		cont.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	if modelInfo.Id == ""{
		statusCode := http.StatusNotFound
		errMessage := fmt.Sprintf("Record not found with id: %s", id)
		logging.ERROR("Record not found, send status code: ", statusCode)
		cont.JSON(statusCode, gin.H{
			"code":    statusCode,
			"message": errMessage,
		})
		return
	}
	cont.JSON(http.StatusOK, modelInfo)
	return
}

/*
Provides the model details by param model name
*/
func (m *MmeApiHandler) GetModelInfoByName(cont *gin.Context) {
	logging.INFO("Get model info by name API ...")
	modelName := cont.Param("modelName")

	bucketObj := m.dbmgr.GetBucketObject(modelName+os.Getenv("INFO_FILE_POSTFIX"), modelName)
	modelInfoListResp := models.ModelInfoResponse{
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

func (m *MmeApiHandler) UpdateModel(c *gin.Context) {
	logging.INFO("Updating model...")
	id := c.Param("id")
	var modelInfo models.ModelInfo

	if err := c.ShouldBindJSON(&modelInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if id != modelInfo.Id {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID in path and body does not match",
		})
		return
	}

	if err := m.iDB.Update(modelInfo); err != nil {
		logging.ERROR(err)
		return
	}

	logging.INFO("model updated")
	c.JSON(http.StatusOK, gin.H{})
}

func (m *MmeApiHandler) DeleteModel(c *gin.Context) {
	logging.INFO("Deleting model...")
	id := c.Param("id")
	if err := m.iDB.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "modelInfo deleted"})
}
