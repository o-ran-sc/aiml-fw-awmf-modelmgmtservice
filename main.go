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
	"database/sql"
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

	// Step 1: Connect to the default 'postgres' database
	defaultDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable",
		configManager.DB.PG_HOST,
		configManager.DB.PG_USER,
		configManager.DB.PG_PASSWORD,
		configManager.DB.PG_PORT,
	)

	defaultDB, err := sql.Open("postgres", defaultDSN)
	if err != nil {
		logging.INFO(fmt.Sprintf("Failed to connect to default database: %v", err))
	}
	defer defaultDB.Close()

	// Step 2: Check if target database exists
	var exists bool
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", configManager.DB.PG_DBNAME)
	err = defaultDB.QueryRow(checkQuery).Scan(&exists)
	if err != nil {
		logging.INFO(fmt.Sprintf("Failed to check database existence: %v", err))
	}

	// Step 3: Create database if missing
	if !exists {
		createQuery := fmt.Sprintf("CREATE DATABASE %s", configManager.DB.PG_DBNAME)
		_, err = defaultDB.Exec(createQuery)
		if err != nil {
			logging.INFO(fmt.Sprintf("Failed to create database %s: %v", configManager.DB.PG_DBNAME, err))
		}
		logging.INFO(fmt.Sprintf("Database '%s' created successfully.", configManager.DB.PG_DBNAME))
	} else {
		logging.INFO(fmt.Sprintf("Database '%s' already exists.", configManager.DB.PG_DBNAME))
	}

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
	err = db.AutoMigrate(
		&models.ModelRelatedInformation{},
		&models.TargetEnvironment{},
	)
	if err != nil {
		logging.ERROR("Failed to migrate database", "error", err)
		os.Exit(-1)
	}

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
