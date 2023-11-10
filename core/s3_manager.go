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
package core

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Manager struct {
	//S3Client has s3 endpoint connection pointer,
	//Which will be used by all s3 bucket related operatios,
	//using fuction to struct binding
	S3Client *s3.S3
}

/*
This function return an pointer to instance of S3_manager struct
the struct instance hold pointer to s3.S3 connection, which is
preconfigured using enviroment variables, such as aws s3
endpoints connection details.
*/
func NewS3Manager() *S3Manager {
	endpoint := os.Getenv("S3_URL")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretAccessKey := os.Getenv("S3_SECRET_KEY")

	config := aws.Config{
		Endpoint:         aws.String(endpoint),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretAccessKey, ""),
		Region:           aws.String(os.Getenv("S3_REGION")),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess, err := session.NewSession(&config)

	if err != nil {
		panic(err)
	}
	S3Client := s3.New(sess)
	return &S3Manager{S3Client} //Return new instance of S3_manager with all config loaded

}

func (s3manager *S3Manager) CreateBucket(bucketName string) {
	_, err := s3manager.S3Client.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(bucketName)})
	if err != nil {
		panic(err)
	}
	println("Bucket created : ", bucketName)
}

func (s3manager *S3Manager) GetBucketObject(objectName string, bucketName string) []byte {

	var response []byte
	getInputs := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}
	result, err := s3manager.S3Client.GetObject(getInputs)
	defer result.Body.Close()

	if err != nil {
		fmt.Println("Error, can't get fetch object..")
		return response
	} else {
		fmt.Println("Successfully retrieved object...")
	}
	//TODO : Error handling
	response, err = io.ReadAll(result.Body)
	return response
}

func (s3manager *S3Manager) DeleteBucket(client *s3.S3, objectName string, bucketName string) {
	success := s3manager.DeleteBucketObject(client, objectName, bucketName)
	if success {
		deleteBucketInput := &s3.DeleteBucketInput{
			Bucket: aws.String(bucketName),
		}
		client.DeleteBucket(deleteBucketInput)
		fmt.Println("Bucket deleted successfully..")
	} else {
		fmt.Println("Failed to delete the Bucket ...")
	}

}

func (s3manager *S3Manager) DeleteBucketObject(client *s3.S3, objectName string, bucketName string) bool {
	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}
	_, err := client.DeleteObject(deleteInput)
	if err != nil {
		fmt.Println("Can not delete the bucket object")
		return false
	}
	fmt.Println("Object deleted successfully..")
	return true
}

func (s3manager *S3Manager) UploadFile(data *bytes.Reader, file_name string, bucketName string) {

	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(file_name),
		Body:   data,
	}
	_, err := s3manager.S3Client.PutObject(params)
	if err != nil {
		fmt.Println("Error in uploading file to bucket ", err)
	}
	fmt.Println("File uploaded to bucket ", bucketName)
}

func (s3manager *S3Manager) ListBucket(client *s3.S3) {
	input := &s3.ListBucketsInput{}
	result, err := client.ListBuckets(input)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
}

// Return list of objects in the buckets
func GetBucketItems(bucketName string) {

}
