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

package models

type Metadata struct {
	Author string `json:"author" validate:"required"`
	Owner  string `json:"owner"`
}

type TargetEnironment struct {
	PlatformName    string `json:"platformName" validate:"required"`
	EnvironmentType string `json:"environmentType" validate:"required"`
	DependencyList  string `json:"dependencyList" validate:"required"`
}

type ModelInformation struct {
	Metadata       Metadata `json:"metadata" gorm:"embedded" validate:"required"`
	InputDataType  string   `json:"inputDataType" validate:"required"`  // this field will be a Comma Separated List
	OutputDataType string   `json:"outputDataType" validate:"required"` // this field will be a Comma Separated List
	// TODO: gorm doesn't support list, need to find the right way
	// TargetEnvironment []TargetEnironment `json:"targetEnvironment" gorm:"embedded"`
}
type ModelID struct {
	ModelName       string `json:"modelName" validate:"required" gorm:"primaryKey"`
	ModelVersion    string `json:"modelVersion" validate:"required" gorm:"primaryKey"`
	ArtifactVersion string `json:"artifactVersion"`
}

type ModelRelatedInformation struct {
	ModelId          ModelID          `json:"modelId,omitempty" validate:"required" gorm:"embedded;primaryKey"`
	Description      string           `json:"description" validate:"required"`
	ModelInformation ModelInformation `json:"modelInformation" validate:"required" gorm:"embedded"`
	ModelLocation    string           `json:"modelLocation"`
}

type ModelInfoResponse struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
