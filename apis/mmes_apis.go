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
	"fmt"
	"io"
	"net/http"
	"os"

	"example.com/mmes/core"
	"github.com/gin-gonic/gin"
)

type ModelInfo struct {
	ModelName string                 `json:"model-name"`
	RAppId    string                 `json:"rapp-id"`
	Metainfo  map[string]interface{} `json:"meta-info"`
}

type MMESApis struct {
}

func init() {
	fmt.Println("Starting api server...")
	router := gin.Default()

	router.POST("/registerModel", RegisterModel)
	router.GET("/getModelInfo", GetModelInfo)
	router.GET("/getModelInfo/:modelName", GetModelInfoByName)
	router.MaxMultipartMemory = 8 << 20 //8 Mb
	router.POST("/uploadModel/:modelName", UploadModel)
	router.GET("/downloadModel/:modelName", DownloadModel)
	router.Run(os.Getenv("MMES_URL"))
	fmt.Println("Started api server...")
}

func RegisterModel(cont *gin.Context) {
	var returnCode int = http.StatusCreated
	var responseMsg string = "Model registered successfully"

	fmt.Println("Creating model...")
	bodyBytes, _ := io.ReadAll(cont.Request.Body)

	var modelInfo ModelInfo
	//Need to unmarshal JSON to Struct, to access request
	//data such as model name, rapp id etc
	err := json.Unmarshal(bodyBytes, &modelInfo)
	if err != nil || modelInfo.ModelName == "" {
		fmt.Println("Error in unmarshalling")
		cont.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": string("Can not parse input data, provide mandatory details"),
		})
		//return
	} else {
		fmt.Println(modelInfo.ModelName, modelInfo.RAppId, modelInfo.Metainfo)
		modelInfoBytes, _ := json.Marshal(modelInfo)

		//TODO Create singleton for s3_manager
		s3_manager := core.NewS3Manager()
		s3Err := s3_manager.CreateBucket(modelInfo.ModelName)
		if s3Err == nil {
			s3_manager.UploadFile(modelInfoBytes, modelInfo.ModelName+os.Getenv("INFO_FILE_PREFIX"), modelInfo.ModelName)
		} else {
			returnCode = http.StatusInternalServerError
			responseMsg = s3Err.Error()
		}
		cont.JSON(returnCode, gin.H{
			"code":    returnCode,
			"message": responseMsg,
		})
	}
}

/*
This API retrieves model info for given model name
input :

	Model name : string
*/
func GetModelInfo(cont *gin.Context) {
	fmt.Println("Fetching model")
	bodyBytes, _ := io.ReadAll(cont.Request.Body)
	//TODO Error checking of request is not in json, i.e. etra ',' at EOF
	jsonMap := make(map[string]interface{})
	json.Unmarshal(bodyBytes, &jsonMap)
	model_name := jsonMap["model-name"].(string)
	fmt.Println("The request model name: ", model_name)

	s3_manager := core.NewS3Manager()
	model_info := s3_manager.GetBucketObject(model_name+os.Getenv("INFO_FILE_PREFIX"), model_name)

	cont.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": string(model_info),
	})
}

/*
Provides the model details by param model name
*/
func GetModelInfoByName(cont *gin.Context) {
	fmt.Println("Get model info by name API ...")
	modelName := cont.Param("modelName")

	s3_manager := core.NewS3Manager()
	model_info := s3_manager.GetBucketObject(modelName+os.Getenv("INFO_FILE_PREFIX"), modelName)

	cont.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": string(model_info),
	})
}

// API to upload the trained model in zip format
// TODO : Model version as input

func UploadModel(cont *gin.Context) {
	fmt.Println("Uploading model API ...")
	modelName := cont.Param("modelName")
	//TODO convert multipart.FileHeader to []byted
	fileHeader, _ := cont.FormFile("file")
	//TODO : Accept only .zip file for trained model
	file, _ := fileHeader.Open()
	//TODO: Handle error response
	defer file.Close()
	byteFile, _ := io.ReadAll((file))

	fmt.Println("Uploading model : ", modelName)
	fmt.Println("Recieved file name :", fileHeader.Filename)

	s3_manager := core.NewS3Manager()
	s3_manager.UploadFile(byteFile, modelName+os.Getenv("MODEL_NAME_PREFIX"), modelName)
	cont.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": string("Model uploaded successfully.."),
	})
}

/*
API to download the trained model from s3 bucket
Input: model name in path params as "modelName"
*/
func DownloadModel(cont *gin.Context) {
	fmt.Println("Download model API ...")
	modelName := cont.Param("modelName")
	fileName := modelName + os.Getenv("MODEL_NAME_PREFIX")
	s3_manager := core.NewS3Manager()
	fileByes := s3_manager.GetBucketObject(fileName, modelName)

	//Return file in api reponse using byte slice
	cont.Header("Content-Disposition", "attachment;"+fileName)
	cont.Header("Content-Type", "application/octet-stream")
	cont.Data(http.StatusOK, "application/octet", fileByes)
}

func GetModel(cont *gin.Context) {
	fmt.Println("Fetching model")
	cont.IndentedJSON(http.StatusOK, " ")
}

func UpdateModel() {
	fmt.Println("Updating model...")
	return
}

func DeleteModel() {
	fmt.Println("Deleting model...")
	return
}
