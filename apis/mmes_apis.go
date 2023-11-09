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
	"bytes"
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

	router.GET("/getModel", GetModel)
	router.POST("/createModel", CreateModel)
	router.Run(os.Getenv("MMES_URL"))
	fmt.Println("Started api server...")
}

func CreateModel(cont *gin.Context) {
	fmt.Println("Creating model...")
	bodyBytes, _ := io.ReadAll(cont.Request.Body)

	var modelInfo ModelInfo
	//Need to unmarshal JSON to Struct, to access request
	//data such as model name, rapp id etc
	err := json.Unmarshal(bodyBytes, &modelInfo)
	if err != nil {
		fmt.Println("Error in unmarshalling")
	}
	fmt.Println(modelInfo.ModelName, modelInfo.RAppId, modelInfo.Metainfo)

	//modelInfo.RAppId = "newRappId-1" Update Unmarshalled struct as per need
	//Need to convert struct to json to create a io.ReadSeeker instance
	//to insert in to a bucket as file/body
	modelInfo_json, err := json.Marshal(modelInfo)
	modelinfo_reader := bytes.NewReader(modelInfo_json) //bytes.Reader is type of io.ReadSeeker

	//TODO Create singleton for s3_manager
	s3_manager := core.NewS3Manager()
	s3_manager.CreateBucket(modelInfo.ModelName)
	s3_manager.UploadFile(modelinfo_reader, modelInfo.ModelName+"_info.json", modelInfo.ModelName)

	cont.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": string("Model details stored sucessfully"),
	})
}

func GetModel(cont *gin.Context) {
	fmt.Println("Fetching model")
	//s3_manager.GetBucketObject(modelInfo.ModelName+"_info.json", modelInfo.ModelName)
	cont.IndentedJSON(http.StatusOK, "Model stored ")
}

func UpdateModel() {
	fmt.Println("Updating model...")
	return
}

func DeleteModel() {
	fmt.Println("Deleting model...")
	return
}
