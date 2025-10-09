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
package apis_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/apis"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/apis_test/mme_mocks"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/logging"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/models"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/routers"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var registerModelBody = `{
	"id" : "id",
    "modelId": {
        "modelName": "model3",
        "modelVersion" : "2"
    },
    "description": "hello world2",
    "modelInformation": {
        "metadata": {
            "author": "someone"
        },
        "inputDataType": "pdcpBytesDl,pdcpBytesUl,kpi",
        "outputDataType": "c, d"
    }
}`

var invalidRegisterModelBody = `{
	"id" : "id",
    "modelId": {
        "modelName": "model3",
        "modelVersion" : "2"
    },
    "description": "hello world2",
    "modelInformation": {
        "metadata": {
            "author": "someone"
        },
        "inputDataType": "pdcpBytesDl,pdcpBytesUl,kpi",
        "outputDataType": "c, d"
    }
`

var invalidRegisterModelBody2 = `{
	"id" : "id",
    "modelId": {
        "modelName": "model3",
        "modelsion" : "2"
    },
    "description": "hello world2",
    "modelInformation": {
        "metadata": {
            "author": "someone"
        },
        "inputDataType": "pdcpBytesDl,pdcpBytesUl,kpi",
        "outputDataType": "c, d"
    }
}`

func TestRegisterModel(t *testing.T) {
	os.Setenv("LOG_FILE_NAME", "testing")
	iDBMockInst := new(mme_mocks.IDBMock)
	iDBMockInst.On("Create", mock.Anything).Return(nil)
	handler := apis.NewMmeApiHandler(nil, iDBMockInst)
	router := routers.InitRouter(handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ai-ml-model-registration/v1/model-registrations", strings.NewReader(registerModelBody))
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
}

func TestRegisterModelFailInvalidJson(t *testing.T) {
	os.Setenv("LOG_FILE_NAME", "testing")
	iDBMockInst := new(mme_mocks.IDBMock)
	handler := apis.NewMmeApiHandler(nil, iDBMockInst)
	router := routers.InitRouter(handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ai-ml-model-registration/v1/model-registrations", strings.NewReader(invalidRegisterModelBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	body, _ := io.ReadAll(w.Body)

	assert.Equal(t, `{"status":400,"title":"Bad Request","detail":"The request json is not correct, unexpected EOF"}`, string(body))
}

func TestRegisterModelFailInvalidRequest(t *testing.T) {
	os.Setenv("LOG_FILE_NAME", "testing")
	iDBMockInst := new(mme_mocks.IDBMock)
	handler := apis.NewMmeApiHandler(nil, iDBMockInst)
	router := routers.InitRouter(handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ai-ml-model-registration/v1/model-registrations", strings.NewReader(invalidRegisterModelBody2))
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	body, _ := io.ReadAll(w.Body)

	assert.Equal(t, `{"status":400,"title":"Bad Request","detail":"The request json is not correct as it can't be validated, Key: 'ModelRelatedInformation.ModelId.ModelVersion' Error:Field validation for 'ModelVersion' failed on the 'required' tag"}`, string(body))
}

func TestRegisterModelFailCreateDuplicateModel(t *testing.T) {
	os.Setenv("LOG_FILE_NAME", "testing")
	iDBMockInst := new(mme_mocks.IDBMock)
	iDBMockInst.On("Create", mock.Anything).Return(&pq.Error{Code: pgerrcode.UniqueViolation})
	handler := apis.NewMmeApiHandler(nil, iDBMockInst)
	router := routers.InitRouter(handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ai-ml-model-registration/v1/model-registrations", strings.NewReader(registerModelBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, 409, w.Code)
	body, _ := io.ReadAll(w.Body)

	assert.Equal(t, "{\"status\":409,\"title\":\"Conflict\",\"detail\":\"model name and version combination already present\"}", string(body))
}

func TestRegisterModelFailCreate(t *testing.T) {
	os.Setenv("LOG_FILE_NAME", "testing")
	iDBMockInst := new(mme_mocks.IDBMock)
	iDBMockInst.On("Create", mock.Anything).Return(&pq.Error{Code: pgerrcode.SQLClientUnableToEstablishSQLConnection})
	handler := apis.NewMmeApiHandler(nil, iDBMockInst)
	router := routers.InitRouter(handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ai-ml-model-registration/v1/model-registrations", strings.NewReader(registerModelBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	body, _ := io.ReadAll(w.Body)

	assert.Equal(t, "{\"status\":500,\"title\":\"Internal Server Error\",\"detail\":\"Database error: pq: \"}", string(body))
}

func TestWhenSuccessGetModelInfoList(t *testing.T) {
	os.Setenv("LOG_FILE_NAME", "testing")

	iDBmockInst := new(mme_mocks.IDBMock)
	iDBmockInst.On("GetAll").Return([]models.ModelRelatedInformation{
		{
			Id: "1234",
			ModelId: models.ModelID{
				ModelName:    "test",
				ModelVersion: "v1.0",
			},
			Description: "this is test modelINfo",
			ModelInformation: models.ModelInformation{
				Metadata: models.Metadata{
					Author: "someone",
				},
				InputDataType:  "pdcpBytesDl,pdcpBytesUl,kpi",
				OutputDataType: "c,d",
			},
		},
	}, nil)

	handler := apis.NewMmeApiHandler(nil, iDBmockInst)
	router := routers.InitRouter(handler)
	responseRecorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/ai-ml-model-discovery/v1/models", nil)
	router.ServeHTTP(responseRecorder, req)

	response := responseRecorder.Result()
	body, _ := io.ReadAll(response.Body)

	var modelInfos []models.ModelRelatedInformation
	logging.INFO("modelinfo", "list:", modelInfos)
	json.Unmarshal(body, &modelInfos)

	assert.Equal(t, 200, responseRecorder.Code)
	assert.Equal(t, 1, len(modelInfos))
}

func TestWhenFailGetModelInfoList(t *testing.T) {
	os.Setenv("LOG_FILE_NAME", "testing")

	iDBmockInst2 := new(mme_mocks.IDBMock)
	iDBmockInst2.On("GetAll").Return([]models.ModelRelatedInformation{}, fmt.Errorf("db not available"))

	handler := apis.NewMmeApiHandler(nil, iDBmockInst2)
	router := routers.InitRouter(handler)
	responseRecorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/ai-ml-model-discovery/v1/models", nil)
	router.ServeHTTP(responseRecorder, req)

	response := responseRecorder.Result()
	body, _ := io.ReadAll(response.Body)

	var modelInfoListResp []models.ModelRelatedInformation
	json.Unmarshal(body, &modelInfoListResp)

	assert.Equal(t, 500, responseRecorder.Code)
}

func TestGetModelInfoParamsInvalid(t *testing.T) {
	os.Setenv("LOG_FILE_NAME", "testing")

	iDBMockInst := new(mme_mocks.IDBMock)
	handler := apis.NewMmeApiHandler(nil, iDBMockInst)
	router := routers.InitRouter(handler)
	responseRecorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/ai-ml-model-discovery/v1/models?model-me=qoe2", nil)
	router.ServeHTTP(responseRecorder, req)

	body, _ := io.ReadAll(responseRecorder.Body)
	fmt.Println(responseRecorder)

	assert.Equal(t, 400, responseRecorder.Code)
	assert.Equal(t, `{"status":400,"title":"Bad Request","detail":"Only allowed params are modelname and modelversion"}`, string(body))
}

func TestGetModelInfoByNameSuccess(t *testing.T) {
	os.Setenv("LOG_FILE_NAME", "testing")

	iDBMockInst := new(mme_mocks.IDBMock)
	iDBMockInst.On("GetModelInfoByName").Return([]models.ModelRelatedInformation{
		{
			Id: "1234",
			ModelId: models.ModelID{
				ModelName:    "test",
				ModelVersion: "v1.0",
			},
			Description: "this is test modelINfo",
			ModelInformation: models.ModelInformation{
				Metadata: models.Metadata{
					Author: "someone",
				},
				InputDataType:  "pdcpBytesDl,pdcpBytesUl,kpi",
				OutputDataType: "c,d",
			},
		},
	}, nil)

	handler := apis.NewMmeApiHandler(nil, iDBMockInst)
	router := routers.InitRouter(handler)
	responseRecorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/ai-ml-model-discovery/v1/models?model-name=qoe1", nil)
	router.ServeHTTP(responseRecorder, req)

	response := responseRecorder.Result()
	body, _ := io.ReadAll(response.Body)

	var modelInfos []models.ModelRelatedInformation
	logging.INFO("modelinfo", "list:", modelInfos)
	json.Unmarshal(body, &modelInfos)

	assert.Equal(t, 200, responseRecorder.Code)
	assert.Equal(t, "1234", modelInfos[0].Id)
}

func TestGetModelInfoByNameFail(t *testing.T) {
	os.Setenv("LOG_FILE_NAME", "testing")

	iDBMockInst := new(mme_mocks.IDBMock)
	iDBMockInst.On("GetModelInfoByName").Return([]models.ModelRelatedInformation{}, fmt.Errorf("db not available"))

	handler := apis.NewMmeApiHandler(nil, iDBMockInst)
	router := routers.InitRouter(handler)
	responseRecorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/ai-ml-model-discovery/v1/models?model-name=qoe1", nil)
	router.ServeHTTP(responseRecorder, req)

	body, _ := io.ReadAll(responseRecorder.Body)

	assert.Equal(t, 500, responseRecorder.Code)
	assert.Equal(t, `{"status":500,"title":"Internal Server Error","detail":"Can't fetch the models due to , db not available"}`, string(body))
}

func TestGetModelInfoByNameAndVersionSuccess(t *testing.T) {
	os.Setenv("LOG_FILE_NAME", "testing")

	iDBMockInst := new(mme_mocks.IDBMock)

	modelInfo := models.ModelRelatedInformation{
		Id: "1234",
		ModelId: models.ModelID{
			ModelName:    "test",
			ModelVersion: "v1.0",
		},
		Description: "this is test modelINfo",
		ModelInformation: models.ModelInformation{
			Metadata: models.Metadata{
				Author: "someone",
			},
			InputDataType:  "pdcpBytesDl,pdcpBytesUl,kpi",
			OutputDataType: "c,d",
		},
	}

	iDBMockInst.On("GetModelInfoByNameAndVer").Return(&modelInfo, nil)

	handler := apis.NewMmeApiHandler(nil, iDBMockInst)
	router := routers.InitRouter(handler)
	responseRecorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/ai-ml-model-discovery/v1/models?model-name=test&model-version=v1.0", nil)
	router.ServeHTTP(responseRecorder, req)

	response := responseRecorder.Result()
	body, _ := io.ReadAll(response.Body)

	var modelInfos []models.ModelRelatedInformation
	logging.INFO("modelinfo", "list:", modelInfos)
	json.Unmarshal(body, &modelInfos)

	assert.Equal(t, 200, responseRecorder.Code)
	assert.Equal(t, "1234", modelInfos[0].Id)
}

func TestGetModelInfoByNameAndVersionFail(t *testing.T) {
	os.Setenv("LOG_FILE_NAME", "testing")

	iDBMockInst := new(mme_mocks.IDBMock)
	modelInfo := models.ModelRelatedInformation{}
	iDBMockInst.On("GetModelInfoByNameAndVer").Return(&modelInfo, fmt.Errorf("db not available"))

	handler := apis.NewMmeApiHandler(nil, iDBMockInst)
	router := routers.InitRouter(handler)
	responseRecorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/ai-ml-model-discovery/v1/models?model-name=test&model-version=v1.0", nil)
	router.ServeHTTP(responseRecorder, req)

	response := responseRecorder.Result()
	body, _ := io.ReadAll(response.Body)

	var modelInfos []models.ModelRelatedInformation
	json.Unmarshal(body, &modelInfos)

	assert.Equal(t, 500, responseRecorder.Code)
	assert.Equal(t, `{"status":500,"title":"Internal Server Error","detail":"Can't fetch all the models due to , db not available"}`, string(body))
}

func TestUploadModelSuccess(t *testing.T) {
	os.Setenv("LOG_FILE_NAME", "testing")
	os.Setenv("MODEL_FILE_POSTFIX", ".zip")
	// Setup Mocks
	iDBMockInst := new(mme_mocks.IDBMock)
	modelName := "test-model"
	modelVersion := "1"
	modelArtifactVersion := "1.0.0"
	modelInfo := models.ModelRelatedInformation{
		ModelId: models.ModelID{
			ModelName:       modelName,
			ModelVersion:    modelVersion,
			ArtifactVersion: modelArtifactVersion,
		},
	}
	iDBMockInst.On("GetModelInfoByNameAndVer").Return(&modelInfo, nil)

	dbMgrMockInst := new(mme_mocks.DbMgrMock)
	dbMgrMockInst.On("UploadFile").Return(nil)
	handler := apis.NewMmeApiHandler(dbMgrMockInst, iDBMockInst)
	router := routers.InitRouter(handler)
	responseRecorder := httptest.NewRecorder()

	// Creating Model.zip for upload
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "Model.zip") // Add file field
	assert.NoError(t, err)
	_, err = part.Write([]byte("fake zip file content"))
	assert.NoError(t, err)
	writer.Close()

	// Upload model
	url := fmt.Sprintf("/ai-ml-model-registration/v1/uploadModel/%s/%s", modelName, modelVersion)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(responseRecorder, req)

	response := responseRecorder.Result()
	responseBody, _ := io.ReadAll(response.Body)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	var responseJson map[string]any
	err = json.Unmarshal(responseBody, &responseJson)
	if err != nil {
		t.Errorf("Error to Unmarshal response-body : Error %s", err.Error())
	}

	newModelInfoStr, err := json.Marshal(responseJson["modelinfo"])
	if err != nil {
		t.Errorf("Error to Marshal model-Info : Error %s", err.Error())
	}

	var newModelInfo models.ModelRelatedInformation
	if err := json.Unmarshal(newModelInfoStr, &newModelInfo); err != nil {
		log.Fatal("unmarshal error:", err)
	}
	assert.Equal(t, newModelInfo.ModelId.ArtifactVersion, "1.1.0")
}

func TestUploadModelFailureModelNotRegistered(t *testing.T) {
	os.Setenv("LOG_FILE_NAME", "testing")
	os.Setenv("MODEL_FILE_POSTFIX", ".zip")
	// Setup Mocks
	iDBMockInst := new(mme_mocks.IDBMock)
	modelName := "test-model"
	modelVersion := "1"
	// Returns Empty model, signifying Model is Not registered
	iDBMockInst.On("GetModelInfoByNameAndVer").Return(&models.ModelRelatedInformation{}, nil)
	handler := apis.NewMmeApiHandler(nil, iDBMockInst)
	router := routers.InitRouter(handler)
	responseRecorder := httptest.NewRecorder()

	// Upload model
	url := fmt.Sprintf("/ai-ml-model-registration/v1/uploadModel/%s/%s", modelName, modelVersion)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	router.ServeHTTP(responseRecorder, req)

	response := responseRecorder.Result()
	responseBody, _ := io.ReadAll(response.Body)
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
	assert.Equal(
		t,
		fmt.Sprintf(`{"status":404,"title":"Not Found","detail":"ModelName: %s and modelVersion: %s is not registered, Kindly register it first!"}`, modelName, modelVersion),
		string(responseBody),
	)
}

func TestUploadModelFailureModelUploadFailure(t *testing.T) {
	os.Setenv("LOG_FILE_NAME", "testing")
	os.Setenv("MODEL_FILE_POSTFIX", ".zip")
	// Setup Mocks
	iDBMockInst := new(mme_mocks.IDBMock)
	modelName := "test-model"
	modelVersion := "1"
	modelArtifactVersion := "1.0.0"
	modelInfo := models.ModelRelatedInformation{
		ModelId: models.ModelID{
			ModelName:       modelName,
			ModelVersion:    modelVersion,
			ArtifactVersion: modelArtifactVersion,
		},
	}
	iDBMockInst.On("GetModelInfoByNameAndVer").Return(&modelInfo, nil)

	dbMgrMockInst := new(mme_mocks.DbMgrMock)
	// Simulate Model-upload-failure
	dbMgrMockInst.On("UploadFile").Return(fmt.Errorf("Unable to upload model"))
	handler := apis.NewMmeApiHandler(dbMgrMockInst, iDBMockInst)
	router := routers.InitRouter(handler)
	responseRecorder := httptest.NewRecorder()

	// Creating Model.zip for upload
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "Model.zip") // Add file field
	assert.NoError(t, err)
	_, err = part.Write([]byte("fake zip file content"))
	assert.NoError(t, err)
	writer.Close()

	// Upload model
	url := fmt.Sprintf("/ai-ml-model-registration/v1/uploadModel/%s/%s", modelName, modelVersion)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(responseRecorder, req)

	response := responseRecorder.Result()
	responseBody, _ := io.ReadAll(response.Body)
	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
	assert.Equal(t, `{"code":500,"message":"Unable to upload model"}`, string(responseBody))
}
