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
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/logging"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/models"
	"gorm.io/gorm"
)

type ModelInfoRepository struct {
	db *gorm.DB
}

func NewModelInfoRepository(db *gorm.DB) *ModelInfoRepository {
	return &ModelInfoRepository{
		db: db,
	}
}

func (repo *ModelInfoRepository) Create(modelInfo models.ModelInfo) error {
	repo.db.Create(modelInfo)
	return nil
}

func (repo *ModelInfoRepository) GetByID(id string) (*models.ModelInfo, error) {
	return nil, nil
}

func (repo *ModelInfoRepository) GetAll() ([]models.ModelInfo, error) {
	var modelInfos []models.ModelInfo
	result := repo.db.Find(&modelInfos)
	if result.Error != nil {
		return nil, result.Error
	}
	return modelInfos, nil
}

func (repo *ModelInfoRepository) Update(modelInfo models.ModelInfo) error {
	if err := repo.db.Save(modelInfo).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ModelInfoRepository) Delete(id string) error {
	logging.INFO("id is:", id)
	if err := repo.db.Delete(&models.ModelInfo{}, "id=?", id).Error; err != nil {
		return err
	}
	return nil
}
func (repo *ModelInfoRepository) GetModelInfoByName(modelName string)([]models.ModelInfo, error){
	var modelInfos []models.ModelInfo
	result := repo.db.Where("model_name = ?", modelName).Find(&modelInfos)
	if result.Error != nil {
		return nil, result.Error
	}
	return modelInfos, nil
}

func (repo *ModelInfoRepository) GetModelInfoByNameAndVer(modelName string, modelVersion string)(*models.ModelInfo, error){
	var modelInfo models.ModelInfo
	result := repo.db.Where("model_name = ? AND model_version = ?", modelName, modelVersion).Find(&modelInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &modelInfo, nil
}
func (repo *ModelInfoRepository) GetModelInfoById(id string)(*models.ModelInfo, error){
	logging.INFO("id is:", id)
	var modelInfo models.ModelInfo
	result := repo.db.Where("id = ?", id).Find(&modelInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &modelInfo, nil
}