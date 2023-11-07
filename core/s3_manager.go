package core

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3_manager struct {
	//Have aws properties, so that they are assiccibe
	// to all the function of this file
}

func (s3manager *S3_manager) getBucketObject(client *s3.S3, objectName string, bucketName string) {

	getInputs := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}
	result, err := client.GetObject(getInputs)
	defer result.Body.Close()

	if err != nil {
		fmt.Println("Error, can't get fetch object..")
		return
	} else {
		fmt.Println("Successfully retrieved object...")
		fmt.Println(result)
	}

	outputFile, error := os.Create("Model-1.zip")
	if error != nil {
		fmt.Println("Could not create local file ")
		return
	}
	defer outputFile.Close()
	//err = os.WriteFile("model-1", []byte(result.Body))
	_, wrerr := io.Copy(outputFile, result.Body)

	if wrerr != nil {
		fmt.Println("Failed to write bucket object in to a localfile")
	} else {
		fmt.Println("Successfully, written bucket object to localfile")
	}

}

func (s3manager *S3_manager) deleteBucket(client *s3.S3, objectName string, bucketName string) {
	success := s3manager.deleteBucketObject(client, objectName, bucketName)
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

func (s3manager *S3_manager) deleteBucketObject(client *s3.S3, objectName string, bucketName string) bool {
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

func (s3manager *S3_manager) uploadFile(client *s3.S3, fileName string, bucketName string) {

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error is opening the file..")
	}
	defer file.Close()

	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	}

	_, err = client.PutObject(params)
	if err != nil {
		fmt.Printf("Error in uploading file to bucket ", err)
	}
	fmt.Printf("File uploaded to bucket ")
}

func (s3manager *S3_manager) listBucket(client *s3.S3) {
	input := &s3.ListBucketsInput{}
	result, err := client.ListBuckets(input)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
}

// Return list of objects in the buckets
func getBucketItems(bucketName string) {

}

func (s3manager *S3_manager) createBucket(client *s3.S3, bucketName string) {
	_, err := client.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(bucketName)})
	if err != nil {
		panic(err)
	}
	println("Bucket created in Leofs with GoLang", bucketName)
}
