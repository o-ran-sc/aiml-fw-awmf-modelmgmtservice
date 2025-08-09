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
	"testing"

	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.Exec("PRAGMA foreign_keys = ON;").Error; err != nil {
		t.Fatalf("enable fk: %v", err)
	}
	if err := db.AutoMigrate(&models.ModelRelatedInformation{}, &models.TargetEnvironment{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	return db
}

func TestCreate_Preload_Update_Delete(t *testing.T) {
	db := openTestDB(t)
	repo := NewModelInfoRepository(db)

	// 1) Create
	mri := models.ModelRelatedInformation{
		ModelId: models.ModelID{
			ModelName:    "resnet",
			ModelVersion: "1.0.0",
		},
		Description: "test model",
		ModelInformation: models.ModelInformation{
			Metadata:       models.Metadata{Author: "alice", Owner: "team-a"},
			InputDataType:  "csv",
			OutputDataType: "json",
		},
		TargetEnvironments: []models.TargetEnvironment{
			{PlatformName: "k8s", EnvironmentType: "prod", DependencyList: "cuda=12.1,torch=2.3"},
			{PlatformName: "ec2", EnvironmentType: "stg", DependencyList: "cpu-only"},
		},
	}

	if err := repo.Create(mri); err != nil {
		t.Fatalf("create: %v", err)
	}

	created, err := repo.GetModelInfoByNameAndVer("resnet", "1.0.0")
	if err != nil {
		t.Fatalf("get after create: %v", err)
	}
	if created.Id == "" {
		t.Fatalf("expected non-empty Id after create")
	}
	if len(created.TargetEnvironments) != 2 {
		t.Fatalf("expected 2 target envs, got %d", len(created.TargetEnvironments))
	}
	mri.Id = created.Id

	// 2) Update
	mri.TargetEnvironments = []models.TargetEnvironment{
		{PlatformName: "k8s", EnvironmentType: "prod", DependencyList: "cuda=12.2,torch=2.4"},
	}
	if err := repo.Update(mri); err != nil {
		t.Fatalf("update: %v", err)
	}

	got2, err := repo.GetModelInfoById(mri.Id)
	if err != nil {
		t.Fatalf("get by id after update: %v", err)
	}
	if len(got2.TargetEnvironments) != 1 {
		t.Fatalf("expected 1 target env after replace, got %d", len(got2.TargetEnvironments))
	}

	// 3) Get by Name/Ver
	got3, err := repo.GetModelInfoByNameAndVer("resnet", "1.0.0")
	if err != nil {
		t.Fatalf("get by name+ver: %v", err)
	}
	if got3.Id != mri.Id {
		t.Fatalf("expected same record by name+ver, got different id")
	}

	// 4) Delete
	rows, err := repo.Delete(mri.Id)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if rows != 1 {
		t.Fatalf("expected 1 row deleted, got %d", rows)
	}

	var cnt int64
	if err := db.Model(&models.TargetEnvironment{}).
		Where("model_related_information_id = ?", mri.Id).
		Count(&cnt).Error; err != nil {
		t.Fatalf("count children: %v", err)
	}
	if cnt != 0 {
		t.Fatalf("expected child rows deleted by CASCADE, remain=%d", cnt)
	}
}
