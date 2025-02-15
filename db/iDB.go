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
package db

import (
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/models"
)

type IDB interface {
	Create(modelInfo models.ModelRelatedInformation) error
	GetByID(id string) (*models.ModelRelatedInformation, error)
	GetAll() ([]models.ModelRelatedInformation, error)
	GetModelInfoByName(modelName string) ([]models.ModelRelatedInformation, error)
	GetModelInfoByNameAndVer(modelName string, modelVersion string) (*models.ModelRelatedInformation, error)
	GetModelInfoById(id string) (*models.ModelRelatedInformation, error)
	Update(modelInfo models.ModelRelatedInformation) error
	Delete(id string) (int64, error)
}
