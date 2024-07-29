# ==================================================================================
#
#       Copyright (c) 2023 Samsung Electronics Co., Ltd. All Rights Reserved.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#          http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#
# ==================================================================================
#Base Image
FROM golang:1.21-alpine AS builder

# location in the container
ENV MME_DIR /home/app/
WORKDIR ${MME_DIR}

# Copy sources into the container
COPY . .
# Install dependencies from go.mod
RUN go mod tidy

RUN LOG_FILE_NAME=testing.log go test ./...

#Build all packages from current dir into bin 
RUN go build -o mme_bin .

FROM alpine:3.20

RUN apk update && apk add bash
RUN apk add --no-cache bash

WORKDIR /app

COPY --from=builder /home/app/mme_bin .

#Start the app
ENTRYPOINT ["/app/mme_bin"]

#Expose the ports for external communication
EXPOSE 8082

