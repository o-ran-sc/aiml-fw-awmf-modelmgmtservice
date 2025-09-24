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
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/db"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/models"
	"github.com/stretchr/testify/mock"
)

type IDBMock struct {
	mock.Mock
	db.IDB
}

func (i *IDBMock) Create(modelInfo models.ModelRelatedInformation) error {
	args := i.Called(modelInfo)
	return args.Error(0)
}

func (i *IDBMock) GetByID(id string) (*models.ModelRelatedInformation, error) {
	return nil, nil
}

func (i *IDBMock) GetAll() ([]models.ModelRelatedInformation, error) {
	args := i.Called()
	if _, ok := args.Get(1).(error); !ok {
		return args.Get(0).([]models.ModelRelatedInformation), nil
	} else {
		var emptyModelInfo []models.ModelRelatedInformation
		return emptyModelInfo, args.Error(1)
	}
}

func (i *IDBMock) Update(modelInfo models.ModelRelatedInformation) error {
	return nil
}

func (i *IDBMock) Delete(id string) (int64, error) {
	return 1, nil
}

func (i *IDBMock) GetModelInfoByName(modelName string) ([]models.ModelRelatedInformation, error) {
	args := i.Called()
	if _, ok := args.Get(1).(error); !ok {
		return args.Get(0).([]models.ModelRelatedInformation), nil
	} else {
		var emptyModelInfo []models.ModelRelatedInformation
		return emptyModelInfo, args.Error(1)
	}
}

func (i *IDBMock) GetModelInfoByNameAndVer(modelName string, modelVersion string) (*models.ModelRelatedInformation, error) {
	args := i.Called()

	if _, ok := args.Get(1).(error); !ok {
		return args.Get(0).(*models.ModelRelatedInformation), nil
	} else {
		var emptyModelInfo *models.ModelRelatedInformation
		return emptyModelInfo, args.Error(1)
	}
}
