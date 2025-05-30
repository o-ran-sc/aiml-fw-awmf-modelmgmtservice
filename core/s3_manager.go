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
	"errors"
	"fmt"
	"io"
	"os"

	"sync"

	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/logging"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var Lock = &sync.Mutex{}
var s3MgrInstance *S3Manager

type S3Manager struct {
	//S3Client has s3 endpoint connection pointer,
	//Which will be used by all s3 bucket related operatios,
	//using fuction to struct binding
	S3Client *s3.S3
}

type DBMgr interface {
	CreateBucket(bucketName string) (err error)
	GetBucketObject(objectName string, bucketName string) (BucketObject, error)
	DeleteBucket(client *s3.S3, objectName string, bucketName string)
	DeleteBucketObject(client *s3.S3, objectName string, bucketName string) bool
	UploadFile(dataBytes []byte, file_name string, bucketName string) error
	ListBucket(bucketObjPostfix string) ([]Bucket, error)
	GetBucketItems(bucketName string)
}

// Singleton for S3Manager
func GetDBManagerInstance() DBMgr {
	Lock.Lock()
	defer Lock.Unlock()

	if s3MgrInstance == nil {
		logging.INFO("Creating single instance for S3Manager")
		s3MgrInstance = newS3Manager()
	} else {
		logging.WARN("S3Manager instance already exists")
	}
	return s3MgrInstance
}

/*
This function return an pointer to instance of S3_manager struct
the struct instance hold pointer to s3.S3 connection, which is
preconfigured using enviroment variables, such as aws s3
endpoints connection details.
*/
func newS3Manager() *S3Manager {
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

// Creates s3 bucket for given bucketName, optionally
// returns named error err
func (s3manager *S3Manager) CreateBucket(bucketName string) (err error) {
	_, s3Err := s3manager.S3Client.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(bucketName)})

	if s3Err != nil {
		logging.ERROR("error", "s3Error:", s3Err)
		//Convert the aws to get the code/error msg for api response
		if aerr, ok := s3Err.(awserr.Error); ok {
			err = errors.New(aerr.Message())
			return
		}
	}
	logging.INFO("Bucket created : " + bucketName)
	return nil
}

// objectName : Name of file/object under given bucket
// bucketName : Name of s3 bucket
// TODO: Return error
func (s3manager *S3Manager) GetBucketObject(objectName string, bucketName string) (BucketObject, error) {

	var response []byte
	getInputs := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}
	result, err := s3manager.S3Client.GetObject(getInputs)
	if err != nil {
		logging.ERROR("Error, can't get fetch object..")
		return response, err
	}
	defer result.Body.Close()
	logging.INFO("Successfully retrieved object...")
	response, err = io.ReadAll(result.Body)
	if err != nil {
		logging.ERROR("Recived error while reading body:", err)
		return nil, err
	}
	return response, nil
}

func (s3manager *S3Manager) DeleteBucket(client *s3.S3, objectName string, bucketName string) {
	success := s3manager.DeleteBucketObject(client, objectName, bucketName)
	if success {
		deleteBucketInput := &s3.DeleteBucketInput{
			Bucket: aws.String(bucketName),
		}
		client.DeleteBucket(deleteBucketInput)
		logging.INFO("Bucket deleted successfully..")
	} else {
		logging.ERROR("Failed to delete the Bucket ...")
	}

}

func (s3manager *S3Manager) DeleteBucketObject(client *s3.S3, objectName string, bucketName string) bool {
	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}
	_, err := client.DeleteObject(deleteInput)
	if err != nil {
		logging.WARN("Can not delete the bucket object")
		return false
	}
	logging.INFO("Object deleted successfully..")
	return true
}

func (s3manager *S3Manager) checkIfBucketExists(bucketName string) (bool, error) {
	_, err := s3manager.S3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		// Check if the error means the bucket doesn't exist
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound":
				return false, nil
			default:
				return false, fmt.Errorf("error checking bucket: %v", aerr)
			}
		}
		return false, fmt.Errorf("error checking bucket: %v", err)
	}
	return true, nil
}

func (s3manager *S3Manager) UploadFile(dataBytes []byte, file_name string, bucketName string) error {

	doesBucketExist, err := s3manager.checkIfBucketExists(bucketName)
	if err != nil {
		logging.DEBUG(fmt.Sprintf("unable to check bucket %s existence, Error : %v", bucketName, err))
		return fmt.Errorf("unable to check bucket %s existence, Error : %v", bucketName, err)
	}
	if !doesBucketExist {
		logging.INFO("Bucket " + bucketName + " Doesn't Exist, Creating One")
		if err := s3manager.CreateBucket(bucketName); err != nil {
			logging.DEBUG(fmt.Sprintf("unable to create bucket for uploading-model, Error : %v", err))
			return fmt.Errorf("unable to create bucket for uploading-model, Error : %v", err)
		}
	}

	dataReader := bytes.NewReader(dataBytes) //bytes.Reader is type of io.ReadSeeker

	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(file_name),
		Body:   dataReader,
	}
	_, err = s3manager.S3Client.PutObject(params)
	if err != nil {
		logging.ERROR("Error in uploading file to bucket ", err)
		return err
	}
	logging.INFO("File uploaded to bucket " + bucketName)
	return nil
}

func (s3manager *S3Manager) ListBucket(bucketObjPostfix string) ([]Bucket, error) {
	input := &s3.ListBucketsInput{}
	listBucketsOutput, err := s3manager.S3Client.ListBuckets(input)

	if err != nil {
		logging.ERROR("Can't get bucket list in s3 ", err)
		return []Bucket{}, err
	}

	BucketList := []Bucket{}
	for _, bucket := range listBucketsOutput.Buckets {
		if bucket.Name == nil {
			continue
		}

		bucketObject, err := s3manager.GetBucketObject(*bucket.Name+bucketObjPostfix, *bucket.Name)
		if err != nil {
			logging.ERROR("Unable to list bucketname ", *bucket.Name, ": Error : ", err.Error())
			continue
		}
		if len(bucketObject) == 0 {
			continue
		}

		BucketList = append(BucketList, Bucket{
			Name:   *bucket.Name,
			Object: bucketObject,
		})
	}

	return BucketList, nil
}

// Return list of objects in the buckets
func (S3Manager *S3Manager) GetBucketItems(bucketName string) {
}
