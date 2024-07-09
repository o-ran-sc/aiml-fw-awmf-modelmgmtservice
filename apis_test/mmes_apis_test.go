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
	"model-name": "test-model",
	"rapp-id": "1234",
	"meta-info": {
		"a": "b"
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
