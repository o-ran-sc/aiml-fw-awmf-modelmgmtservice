/*
==================================================================================
Copyright (c) 2025 Minje Kim <alswp006@gmail.com> All Rights Reserved.

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

	d, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := d.AutoMigrate(
		&models.ModelRelatedInformation{},
		&models.TargetEnvironment{},
	); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	return d
}

func newRepo(t *testing.T) *ModelInfoRepository {
	return NewModelInfoRepository(openTestDB(t))
}

func mkMRI(name, ver string, envs []models.TargetEnvironment) models.ModelRelatedInformation {
	return models.ModelRelatedInformation{
		ModelId:     models.ModelID{ModelName: name, ModelVersion: ver},
		Description: "test",
		ModelInformation: models.ModelInformation{
			Metadata:          models.Metadata{Author: "tester"},
			InputDataType:     "csv",
			OutputDataType:    "json",
			TargetEnvironment: envs,
		},
	}
}
func TestCRUD_Subtests(t *testing.T) {
	repo := newRepo(t)

	var id string

	t.Run("Create_with_envs", func(t *testing.T) {
		m := mkMRI("resnet", "1.0.0", []models.TargetEnvironment{
			{PlatformName: "k8s", EnvironmentType: "prod", DependencyList: "cuda=12.1,torch=2.3"},
			{PlatformName: "ec2", EnvironmentType: "stg", DependencyList: "cpu-only"},
		})
		if err := repo.Create(m); err != nil {
			t.Fatalf("create: %v", err)
		}
		got, err := repo.GetModelInfoByNameAndVer("resnet", "1.0.0")
		if err != nil {
			t.Fatalf("get after create: %v", err)
		}
		if got.Id == "" {
			t.Fatalf("expected non-empty id")
		}
		if len(got.ModelInformation.TargetEnvironment) != 2 {
			t.Fatalf("want 2 envs, got %d", len(got.ModelInformation.TargetEnvironment))
		}
		id = got.Id
	})

	t.Run("Read", func(t *testing.T) {
		byID, err := repo.GetModelInfoById(id)
		if err != nil {
			t.Fatalf("get by id: %v", err)
		}
		if byID.Id != id {
			t.Fatalf("id mismatch")
		}

		all, err := repo.GetAll()
		if err != nil {
			t.Fatalf("get all: %v", err)
		}
		if len(all) == 0 {
			t.Fatalf("expected >=1")
		}

		list, err := repo.GetModelInfoByName("resnet")
		if err != nil {
			t.Fatalf("get by name: %v", err)
		}
		if len(list) == 0 {
			t.Fatalf("expected >=1 resnet")
		}
	})

	t.Run("Update_Replace_with_new_values", func(t *testing.T) {
		cur, err := repo.GetModelInfoById(id)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		cur.ModelInformation.TargetEnvironment = []models.TargetEnvironment{
			{PlatformName: "edge", EnvironmentType: "prod", DependencyList: "cuda=12.2,torch=2.4"},
		}
		if err := repo.Update(*cur); err != nil {
			t.Fatalf("update replace: %v", err)
		}
		after, err := repo.GetModelInfoById(id)
		if err != nil {
			t.Fatalf("get after replace: %v", err)
		}
		if len(after.ModelInformation.TargetEnvironment) != 1 {
			t.Fatalf("want 1 env after replace, got %d", len(after.ModelInformation.TargetEnvironment))
		}
		if after.ModelInformation.TargetEnvironment[0].PlatformName != "edge" {
			t.Fatalf("unexpected platform after replace: %+v", after.ModelInformation.TargetEnvironment[0])
		}
	})

	t.Run("Update_Partial_nil_keeps_existing", func(t *testing.T) {
		cur, err := repo.GetModelInfoById(id)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		cur.ModelInformation.TargetEnvironment = nil
		if err := repo.Update(*cur); err != nil {
			t.Fatalf("update partial nil: %v", err)
		}
		after, err := repo.GetModelInfoById(id)
		if err != nil {
			t.Fatalf("get after partial: %v", err)
		}
		if len(after.ModelInformation.TargetEnvironment) != 1 {
			t.Fatalf("want keep 1 env, got %d", len(after.ModelInformation.TargetEnvironment))
		}
	})

	t.Run("Update_Clear_with_empty_slice", func(t *testing.T) {
		cur, err := repo.GetModelInfoById(id)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		cur.ModelInformation.TargetEnvironment = []models.TargetEnvironment{}
		if err := repo.Update(*cur); err != nil {
			t.Fatalf("update clear []: %v", err)
		}
		after, err := repo.GetModelInfoById(id)
		if err != nil {
			t.Fatalf("get after clear: %v", err)
		}
		if len(after.ModelInformation.TargetEnvironment) != 0 {
			t.Fatalf("want 0 env after clear, got %d", len(after.ModelInformation.TargetEnvironment))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		rows, err := repo.Delete(id)
		if err != nil {
			t.Fatalf("delete: %v", err)
		}
		if rows != 1 {
			t.Fatalf("want 1 row deleted, got %d", rows)
		}
		if _, err := repo.GetModelInfoById(id); err == nil {
			t.Fatalf("expected error after delete, got nil")
		}
	})
}
