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
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/apis"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/core"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/routers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var registerModelBody = `{
	"id" : "id", 
	"model-id": {
		"modelName": "test-model",
		"modelVersion":"1"
	},
	"description": "testing",
	"meta-info": {
		"metadata": {
			"author":"tester"
		}
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

func TestRegisterModel(t *testing.T) {
	os.Setenv("LOG_FILE_NAME", "testing")
	dbMgrMockInst := new(dbMgrMock)
	dbMgrMockInst.On("CreateBucket", "test-model").Return(nil)
	handler := apis.NewMmeApiHandler(dbMgrMockInst)
	router := routers.InitRouter(handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/registerModel", strings.NewReader(registerModelBody))
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
}

func TestWhenSuccessGetModelInfoList(t *testing.T) {
	// Setting ENV
	os.Setenv("LOG_FILE_NAME", "testing")

	// Setting Mock
	dbMgrMockInst := new(dbMgrMock)
	dbMgrMockInst.On("ListBucket").Return([]core.Bucket{
		{
			Name:   "qoe",
			Object: []byte(registerModelBody),
		},
	}, nil)

	handler := apis.NewMmeApiHandler(dbMgrMockInst)
	router := routers.InitRouter(handler)
	responseRecorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/getModelInfo", nil)
	router.ServeHTTP(responseRecorder, req)

	response := responseRecorder.Result()
	body, _ := io.ReadAll(response.Body)

	var modelInfoListResp struct {
		Code    int                      `json:"code"`
		Message []apis.ModelInfoResponse `json:"message"`
	}
	json.Unmarshal(body, &modelInfoListResp)

	assert.Equal(t, 200, responseRecorder.Code)
	assert.Equal(t, 200, modelInfoListResp.Code)
	assert.Equal(t, registerModelBody, modelInfoListResp.Message[0].Data)
}

func TestWhenFailGetModelInfoList(t *testing.T) {
	// Setting ENV
	os.Setenv("LOG_FILE_NAME", "testing")

	// Setting Mock
	dbMgrMockInst := new(dbMgrMock)
	dbMgrMockInst.On("ListBucket").Return([]core.Bucket{}, errors.New("Test: Fail GetModelInfoList"))

	handler := apis.NewMmeApiHandler(dbMgrMockInst)
	router := routers.InitRouter(handler)
	responseRecorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/getModelInfo", nil)
	router.ServeHTTP(responseRecorder, req)

	response := responseRecorder.Result()
	body, _ := io.ReadAll(response.Body)

	var modelInfoListResp struct {
		Code    int                      `json:"code"`
		Message []apis.ModelInfoResponse `json:"message"`
	}
	json.Unmarshal(body, &modelInfoListResp)

	assert.Equal(t, 500, responseRecorder.Code)
	assert.Equal(t, 500, modelInfoListResp.Code)
	assert.Equal(t, 0, len(modelInfoListResp.Message))
}
