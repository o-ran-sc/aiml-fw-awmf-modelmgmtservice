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

type S3_manager struct {
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
func New_s3manager() *S3_manager {
	endpoint := os.Getenv("S3_URL")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	SecretAccessKey := os.Getenv("S3_SECRETE_KEY")

	config := aws.Config{
		Endpoint:         aws.String(endpoint),
		Credentials:      credentials.NewStaticCredentials(accessKey, SecretAccessKey, ""),
		Region:           aws.String(os.Getenv("S3_REGION")),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess, err := session.NewSession(&config)

	if err != nil {
		panic(err)
	}
	S3Client := s3.New(sess)
	return &S3_manager{S3Client} //Return new instance of S3_manager with all config loaded

}

func (s3manager *S3_manager) CreateBucket(bucketName string) {
	_, err := s3manager.S3Client.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(bucketName)})
	if err != nil {
		panic(err)
	}
	println("Bucket created : ", bucketName)
}

func (s3manager *S3_manager) GetBucketObject(objectName string, bucketName string) {

	getInputs := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}
	result, err := s3manager.S3Client.GetObject(getInputs)
	defer result.Body.Close()

	if err != nil {
		fmt.Println("Error, can't get fetch object..")
		return
	} else {
		fmt.Println("Successfully retrieved object...")
	}

	outputFile, error := os.Create(bucketName + ".json")
	if error != nil {
		fmt.Println("Could not create local file ")
		return
	}
	defer outputFile.Close()
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

func (s3manager *S3_manager) UploadFile(data *bytes.Reader, file_name string, bucketName string) {

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
