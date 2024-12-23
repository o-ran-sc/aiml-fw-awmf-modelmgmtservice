/*
==================================================================================
Copyright (c) 2023 Samsung Electronics Co., Ltd. All Rights Reserved.

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
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/apis"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/config"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/core"
	modelDB "gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/db"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/logging"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/models"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/routers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := config.Load(config.NewConfigDataValidator(), config.NewEnvDataLoader(nil)); err != nil {
		logging.ERROR("error in loading config", "error", err)
		os.Exit(-1)
	}

	configManager := config.GetConfigManager()
	logging.INFO("config mgr prepared", "configmgr", configManager)
	// setup the database connection

	// connection string
	DSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable ",
		configManager.DB.PG_HOST,
		configManager.DB.PG_USER,
		configManager.DB.PG_PASSWORD,
		configManager.DB.PG_DBNAME,
		configManager.DB.PG_PORT,
	)

	logging.INFO("Connection string for DB", DSN)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  DSN,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	// db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		logging.ERROR("database not available")
		os.Exit(-1)
	}

	// Auto migrate the scheme
	db.AutoMigrate(&models.ModelRelatedInformation{})
	repo := modelDB.NewModelInfoRepository(db)

	router := routers.InitRouter(
		apis.NewMmeApiHandler(
			core.GetDBManagerInstance(),
			repo,
		))
	server := http.Server{
		Addr:         configManager.App.MMES_URL,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logging.INFO("Starting api..")
	err = server.ListenAndServe()
	if err != nil {
		logging.ERROR("error", err)
	}
}
