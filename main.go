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
	"net/http"
	"os"
	"time"

	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/apis"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/core"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/logging"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/routers"
)

func main() {
	router := routers.InitRouter(
		apis.NewMmeApiHandler(
			core.GetDBManagerInstance(),
		))
	server := http.Server{
		Addr:         os.Getenv("MMES_URL"),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logging.INFO("Starting api..")
	server.ListenAndServe()
}
