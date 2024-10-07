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
	Author string `json:"author"`
}

type ModelSpec struct {
	Metadata Metadata `json:"metadata" gorm:"embedded"`
}
type ModelID struct {
	ModelName    string `json:"modelName"`
	ModelVersion string `json:"modelVersion"`
}

type ModelInfo struct {
	Id          string    `json:"id" gorm:"primaryKey"`
	ModelId     ModelID   `json:"model-id,omitempty" gorm:"embedded"`
	Description string    `json:"description"`
	ModelSpec   ModelSpec `json:"meta-info" gorm:"embedded"`
}

type ModelInfoResponse struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
