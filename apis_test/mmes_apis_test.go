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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/apis"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/core"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/db"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/logging"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/models"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/routers"
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

type dbMgrMock struct {
	mock.Mock
	core.DBMgr
}

func (d *dbMgrMock) CreateBucket(bucketName string) (err error) {
	args := d.Called(bucketName)
	return args.Error(0)
}

func (d *dbMgrMock) UploadFile(dataBytes []byte, file_name string, bucketName string) {
}

func (d *dbMgrMock) ListBucket(bucketObjPostfix string) ([]core.Bucket, error) {
	args := d.Called()
	return args.Get(0).([]core.Bucket), args.Error(1)
}

type iDBMock struct {
	mock.Mock
	db.IDB
}

func (i *iDBMock) Create(modelInfo models.ModelRelatedInformation) error {
	args := i.Called(modelInfo)
	return args.Error(0)
}
func (i *iDBMock) GetByID(id string) (*models.ModelRelatedInformation, error) {
	return nil, nil
}
func (i *iDBMock) GetAll() ([]models.ModelRelatedInformation, error) {
	args := i.Called()
	if _, ok := args.Get(1).(error); !ok {
		return args.Get(0).([]models.ModelRelatedInformation), nil
	} else {
		var emptyModelInfo []models.ModelRelatedInformation
		return emptyModelInfo, args.Error(1)
	}
}
func (i *iDBMock) Update(modelInfo models.ModelRelatedInformation) error {
	return nil
}
func (i *iDBMock) Delete(id string) (int64, error) {
	return 1, nil
}

func TestRegisterModel(t *testing.T) {
	os.Setenv("LOG_FILE_NAME", "testing")
	iDBMockInst := new(iDBMock)
	iDBMockInst.On("Create", mock.Anything).Return(nil)
	handler := apis.NewMmeApiHandler(nil, iDBMockInst)
	router := routers.InitRouter(handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ai-ml-model-registration/v1/model-registrations", strings.NewReader(registerModelBody))
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
}

func TestWhenSuccessGetModelInfoList(t *testing.T) {
	// Setting ENV
	os.Setenv("LOG_FILE_NAME", "testing")

	// Setting Mock
	iDBmockInst := new(iDBMock)
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
	// Setting ENV
	os.Setenv("LOG_FILE_NAME", "testing")

	// Setting Mock
	iDBmockInst2 := new(iDBMock)
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
