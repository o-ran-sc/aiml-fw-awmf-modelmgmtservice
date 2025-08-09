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

func (repo *ModelInfoRepository) Create(modelInfo models.ModelRelatedInformation) error {
	result := repo.db.Create(&modelInfo)
	return result.Error
}

func (repo *ModelInfoRepository) GetByID(id string) (*models.ModelRelatedInformation, error) {
	return nil, nil
}

func (repo *ModelInfoRepository) GetAll() ([]models.ModelRelatedInformation, error) {
	var modelInfos []models.ModelRelatedInformation
	result := repo.db.Session(&gorm.Session{SkipHooks: true}).Find(&modelInfos)
	if result.Error != nil {
		return nil, result.Error
	}
	if err := attachEnvsBatch(repo.db, modelInfos); err != nil {
		return nil, err
	}
	return modelInfos, nil
}

func (repo *ModelInfoRepository) Update(modelInfo models.ModelRelatedInformation) error {
	if err := repo.db.Save(&modelInfo).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ModelInfoRepository) Delete(id string) (int64, error) {
	result := repo.db.Delete(&models.ModelRelatedInformation{}, "id = ?", id)
	return result.RowsAffected, result.Error
}

func (repo *ModelInfoRepository) GetModelInfoByName(modelName string) ([]models.ModelRelatedInformation, error) {
	var modelInfos []models.ModelRelatedInformation
	result := repo.db.Session(&gorm.Session{SkipHooks: true}).Where("model_name = ?", modelName).Find(&modelInfos)
	if result.Error != nil {
		return nil, result.Error
	}
	if err := attachEnvsBatch(repo.db, modelInfos); err != nil {
		return nil, err
	}
	return modelInfos, nil
}

func (repo *ModelInfoRepository) GetModelInfoByNameAndVer(modelName string, modelVersion string) (*models.ModelRelatedInformation, error) {
	var m models.ModelRelatedInformation
	result := repo.db.Session(&gorm.Session{SkipHooks: true}).
		Where("model_name = ? AND model_version = ?", modelName, modelVersion).
		First(&m)
	if result.Error != nil {
		return nil, result.Error
	}
	if err := attachEnvsOne(repo.db, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (repo *ModelInfoRepository) GetModelInfoById(id string) (*models.ModelRelatedInformation, error) {
	logging.INFO("id is:", id)
	var m models.ModelRelatedInformation
	result := repo.db.Session(&gorm.Session{SkipHooks: true}).
		Where("id = ?", id).
		First(&m)
	if result.Error != nil {
		return nil, result.Error
	}
	if err := attachEnvsOne(repo.db, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func attachEnvsOne(tx *gorm.DB, m *models.ModelRelatedInformation) error {
	if m == nil || m.Id == "" {
		return nil
	}
	var rows []struct {
		PlatformName    string
		EnvironmentType string
		DependencyList  string
	}
	if err := tx.Table("target_environments").
		Select("platform_name, environment_type, dependency_list").
		Where("model_related_information_id = ?", m.Id).
		Find(&rows).Error; err != nil {
		return err
	}
	m.ModelInformation.TargetEnvironment = m.ModelInformation.TargetEnvironment[:0]
	for _, r := range rows {
		m.ModelInformation.TargetEnvironment = append(
			m.ModelInformation.TargetEnvironment,
			models.TargetEnvironment{
				PlatformName:    r.PlatformName,
				EnvironmentType: r.EnvironmentType,
				DependencyList:  r.DependencyList,
			},
		)
	}
	return nil
}

func attachEnvsBatch(tx *gorm.DB, parents []models.ModelRelatedInformation) error {
	if len(parents) == 0 {
		return nil
	}

	ids := make([]string, 0, len(parents))
	pos := make(map[string]int, len(parents))
	for i := range parents {
		if id := parents[i].Id; id != "" {
			ids = append(ids, id)
			pos[id] = i
		}
	}
	if len(ids) == 0 {
		return nil
	}

	var rows []struct {
		ModelRelatedInformationID string
		PlatformName              string
		EnvironmentType           string
		DependencyList            string
	}
	if err := tx.Table("target_environments").
		Select("model_related_information_id, platform_name, environment_type, dependency_list").
		Where("model_related_information_id IN ?", ids).
		Find(&rows).Error; err != nil {
		return err
	}

	for _, r := range rows {
		if i, ok := pos[r.ModelRelatedInformationID]; ok {
			parents[i].ModelInformation.TargetEnvironment = append(
				parents[i].ModelInformation.TargetEnvironment,
				models.TargetEnvironment{
					PlatformName:    r.PlatformName,
					EnvironmentType: r.EnvironmentType,
					DependencyList:  r.DependencyList,
				},
			)
		}
	}
	return nil
}
