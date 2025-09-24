/*
==================================================================================
Copyright (c) 2025 Samsung Electronics Co., Ltd. All Rights Reserved.

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
package mme_mocks

import (
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/core"
	"github.com/stretchr/testify/mock"
)

type DbMgrMock struct {
	mock.Mock
	core.DBMgr
}

func (d *DbMgrMock) CreateBucket(bucketName string) (err error) {
	args := d.Called(bucketName)
	return args.Error(0)
}

func (d *DbMgrMock) UploadFile(dataBytes []byte, file_name string, bucketName string) error {
	args := d.Called()
	// If error is passed, return the error
	if _, ok := args.Get(0).(error); ok {
		return args.Get(0).(error)
	}

	return nil
}

func (d *DbMgrMock) ListBucket(bucketObjPostfix string) ([]core.Bucket, error) {
	args := d.Called()
	return args.Get(0).([]core.Bucket), args.Error(1)
}
