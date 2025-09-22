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
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/core"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/db"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/logging"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

const (
	MODELNAME    = "model-name"
	MODELVERSION = "model-version"
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
	logging.INFO("registering model info")

	var modelInfo models.ModelRelatedInformation

	if err := cont.ShouldBindJSON(&modelInfo); err != nil {
		cont.JSON(http.StatusBadRequest, models.ProblemDetail{
			Status: http.StatusBadRequest,
			Title:  "Bad Request",
			Detail: fmt.Sprintf("The request json is not correct, %s", err.Error()),
		})
		return
	}

	id := uuid.New()
	modelInfo.Id = id.String()

	validate := validator.New()
	if err := validate.Struct(modelInfo); err != nil {
		cont.JSON(http.StatusBadRequest, models.ProblemDetail{
			Status: http.StatusBadRequest,
			Title:  "Bad Request",
			Detail: fmt.Sprintf("The request json is not correct as it can't be validated, %s", err.Error()),
		})
		return
	}

	// by default when a model is registered its artifact version is set to 0.0.0
	modelInfo.ModelId.ArtifactVersion = "0.0.0"

	if err := m.iDB.Create(modelInfo); err != nil {
		logging.ERROR("error", err)
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == pgerrcode.UniqueViolation {
				cont.JSON(http.StatusConflict, models.ProblemDetail{
					Status: http.StatusConflict,
					Title:  "Conflict",
					Detail: "model name and version combination already present",
				})
				return
			}
			cont.JSON(http.StatusInternalServerError, models.ProblemDetail{
				Status: http.StatusInternalServerError,
				Title:  "Internal Server Error",
				Detail: fmt.Sprintf("Database error: %s", err.Error()),
			})
			return
		}
	}

	logging.INFO("model is saved.")
	cont.Header("Location", "ai-ml-model-registration/v1/model-registrations/"+id.String())
	cont.JSON(http.StatusCreated, gin.H{
		"modelInfo": modelInfo,
	})
}

/*
This API retrieves model info list managed in modelmgmtservice
*/
func (m *MmeApiHandler) GetModelInfo(cont *gin.Context) {
	logging.INFO("Get model info ")
	queryParams := cont.Request.URL.Query()
	// to check only modelName and modelVersion can be passed.
	allowedParams := map[string]bool{
		MODELNAME:    true,
		MODELVERSION: true,
	}

	for key := range queryParams {
		if !allowedParams[key] {
			logging.ERROR("error:", "Only allowed params are modelname and modelversion")
			cont.JSON(http.StatusBadRequest, models.ProblemDetail{
				Status: http.StatusBadRequest,
				Title:  "Bad Request",
				Detail: fmt.Sprintf("Only allowed params are modelname and modelversion"),
			})
			return
		}
	}

	modelName := cont.Query(MODELNAME)
	modelVersion := cont.Query(MODELVERSION)

	if modelName == "" && modelVersion == "" {
		// return all modelinfo stored
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
			modelInfos, err := m.iDB.GetModelInfoByName(modelName)
			if err != nil {
				statusCode := http.StatusInternalServerError
				logging.ERROR("Error occurred, send status code: ", statusCode)
				cont.JSON(statusCode, models.ProblemDetail{
					Status: http.StatusInternalServerError,
					Title:  "Internal Server Error",
					Detail: fmt.Sprintf("Can't fetch the models due to , %s", err.Error()),
				})
				return
			}
			cont.JSON(http.StatusOK, modelInfos)
			return
		} else {
			// get all modelInfo by model name and version
			modelInfo, err := m.iDB.GetModelInfoByNameAndVer(modelName, modelVersion)
			if err != nil {
				statusCode := http.StatusInternalServerError
				logging.ERROR("Error occurred, send status code: ", statusCode)
				cont.JSON(statusCode, models.ProblemDetail{
					Status: http.StatusInternalServerError,
					Title:  "Internal Server Error",
					Detail: fmt.Sprintf("Can't fetch all the models due to , %s", err.Error()),
				})
				return
			}
			if modelInfo.ModelId.ModelName != modelName && modelInfo.ModelId.ModelVersion != modelVersion {
				statusCode := http.StatusNotFound
				logging.ERROR("Record not found, send status code: ", statusCode)
				cont.JSON(statusCode, models.ProblemDetail{
					Status: http.StatusNotFound,
					Title:  "Not Found",
					Detail: fmt.Sprintf("Record not found with modelName: %s and modelVersion: %s", modelName, modelVersion),
				})
				return
			}
			response := []models.ModelRelatedInformation{*modelInfo}
			cont.JSON(http.StatusOK, response)
			return
		}
	}
}

func (m *MmeApiHandler) GetModelInfoById(cont *gin.Context) {
	logging.INFO("Get model info by id ...")
	id := cont.Param("modelRegistrationId")
	modelInfo, err := m.iDB.GetModelInfoById(id)
	if err != nil {
		logging.ERROR("error:", err)
		cont.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	if modelInfo.Id == "" {
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

	bucketObj, err := m.dbmgr.GetBucketObject(modelName+os.Getenv("INFO_FILE_POSTFIX"), modelName)
	if err != nil {
		logging.ERROR("Unable to GetModelInfoByName: Error ", err)
		cont.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	modelInfoListResp := models.ModelInfoResponse{
		Name: modelName,
		Data: string(bucketObj),
	}

	cont.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": modelInfoListResp,
	})
}

/*
* The following API uploads the trained model in zip-format to the provided modelId (modelName, modelVersion, artifactVersion)
* Note: Model MUST be registered first, and then should be uploaded
 */
func (m *MmeApiHandler) UploadModel(cont *gin.Context) {
	logging.INFO("Uploading model API ...")
	modelName := cont.Param("modelName")
	modelVersion := cont.Param("modelVersion")
	artifactVersion := cont.Param("artifactVersion")

	// Confirm if Model with Given ModelId: (ModelName and ModelVersion) is Registered or not:
	modelInfo, err := m.iDB.GetModelInfoByNameAndVer(modelName, modelVersion)
	if err != nil {
		statusCode := http.StatusInternalServerError
		logging.ERROR("Error occurred while getting models: " + strconv.Itoa(statusCode))
		cont.JSON(statusCode, models.ProblemDetail{
			Status: statusCode,
			Title:  "Internal Server Error",
			Detail: fmt.Sprintf("Can't fetch model with modelName : %s & modelVersion : %s due to , %s", modelName, modelVersion, err.Error()),
		})
		return
	}

	if modelInfo.ModelId.ModelName != modelName && modelInfo.ModelId.ModelVersion != modelVersion {
		statusCode := http.StatusNotFound
		logging.ERROR("Record not found, send status code: " + strconv.Itoa(statusCode))
		cont.JSON(statusCode, models.ProblemDetail{
			Status: statusCode,
			Title:  "Model not registered",
			Detail: fmt.Sprintf("ModelName: %s and modelVersion: %s is not registered, Kindly register it first!", modelName, modelVersion),
		})
		return
	}

	modelKey := fmt.Sprintf("%s_%s_%s", modelName, modelVersion, artifactVersion)
	exportBucket := strings.ToLower(modelName)
	//TODO convert multipart.FileHeader to []byte
	fileHeader, _ := cont.FormFile("file")
	//TODO : Accept only .zip file for trained model
	file, _ := fileHeader.Open()
	defer file.Close()
	byteFile, _ := io.ReadAll((file))

	logging.INFO("Uploading model : " + modelKey)
	if err := m.dbmgr.UploadFile(byteFile, modelKey+os.Getenv("MODEL_FILE_POSTFIX"), exportBucket); err != nil {
		logging.ERROR("Failed to Upload Model : ", err)
		cont.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
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
	modelVersion := cont.Param("modelVersion")
	artifactVersion := cont.Param("artifactVersion")

	modelKey := fmt.Sprintf("%s_%s_%s", modelName, modelVersion, artifactVersion)
	exportBucket := strings.ToLower(modelName)

	fileName := modelKey + os.Getenv("MODEL_FILE_POSTFIX")
	fileByes, err := m.dbmgr.GetBucketObject(fileName, exportBucket)
	if err != nil {
		cont.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	// Return file in api response using byte slice
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
	id := c.Param("modelRegistrationId")
	var modelInfo models.ModelRelatedInformation

	if err := c.ShouldBindJSON(&modelInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	existingModelInfo, err := m.iDB.GetModelInfoById(id)

	if err != nil || existingModelInfo.Id == "" {
		statusCode := http.StatusNotFound
		logging.ERROR("Error occurred, send status code: ", statusCode)
		c.JSON(statusCode, gin.H{
			"code":    statusCode,
			"message": fmt.Sprintf("model not found with id: %s", id),
		})
		return
	}

	if existingModelInfo.ModelId.ModelName != modelInfo.ModelId.ModelName || existingModelInfo.ModelId.ModelVersion != modelInfo.ModelId.ModelVersion {
		statusCode := http.StatusBadRequest
		logging.ERROR("Error occurred, send status code: ", statusCode)
		c.JSON(statusCode, gin.H{
			"code":    statusCode,
			"message": fmt.Sprintf("model with id: %s has different modelName and modelVersion than provided", id),
		})
		return
	}

	modelInfo.Id = id
	if err := m.iDB.Update(modelInfo); err != nil {
		logging.ERROR("error in update db", "Error:", err)
		return
	}

	logging.INFO("model updated")
	c.JSON(http.StatusOK, gin.H{
		"modelinfo": modelInfo,
	})
}

func (m *MmeApiHandler) DeleteModel(cont *gin.Context) {
	id := cont.Param("modelRegistrationId")
	logging.INFO("Deleting model... id = ", id)
	_, err := m.iDB.Delete(id)
	if err != nil {
		cont.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	cont.JSON(http.StatusNoContent, nil)
}

func (m *MmeApiHandler) UpdateArtifact(cont *gin.Context) {
	logging.INFO("Update artifact version of model")
	modelname := cont.Param("modelname")
	modelversion := cont.Param("modelversion")
	artifactversion := cont.Param("artifactversion")
	modelInfo, err := m.iDB.GetModelInfoByNameAndVer(modelname, modelversion)
	if err != nil {
		logging.ERROR("error:", err)
		cont.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	modelInfo.ModelId.ArtifactVersion = artifactversion
	if err := m.iDB.Update(*modelInfo); err != nil {
		logging.ERROR("error in update db", "error:", err)
		return
	}
	logging.INFO("model updated")
	cont.JSON(http.StatusOK, gin.H{
		"modelinfo": modelInfo,
	})
}
