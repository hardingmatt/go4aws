package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var bucketName = "mhhardin-bucket-98765"
var keyName = "keyName"

var service *s3.S3

func main() {
	service = setupService()

	createBucket()
	uploadFile()
	readFile()
	deleteFile()
	deleteBucket()
}

func setupService() *s3.S3 {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return nil
	}

	svc := s3.New(sess, aws.NewConfig().WithRegion("us-west-2"))
	return svc
}

func createHeadBucketInput() *s3.HeadBucketInput {
	headBucketInput := &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	}
	return headBucketInput
}

func createBucket() {
	fmt.Println("createBucket")

	if bucketExists() {
		return
	}
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}

	resp, err := service.CreateBucket(input)
	fmt.Println(resp, err)

	service.WaitUntilBucketExists(createHeadBucketInput())
}

func bucketExists() bool {
	fmt.Println("bucketExists")

	_, err := service.HeadBucket(createHeadBucketInput())

	return err == nil
}

func deleteBucket() {
	fmt.Println("deleteBucket")

	input := &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	}
	resp, err := service.DeleteBucket(input)
	fmt.Println(resp, err)

	service.WaitUntilBucketNotExists(createHeadBucketInput())
}

func uploadFile() {
	fmt.Println("uploadFile")

	d := []byte("line 1\nline 2\n")
	ioutil.WriteFile("tmp", d, 0600)
	f, _ := os.Open("tmp")
	defer f.Close()

	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
		Body:   f,
	}
	service.PutObject(input)
	os.Remove("tmp")
}

func readFile() {
	fmt.Println("readFile")

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	}
	resp, err := service.GetObject(input)
	fmt.Println(resp, err)
}

func deleteFile() {
	fmt.Println("deleteFile")

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	}
	resp, err := service.DeleteObject(input)
	fmt.Println(resp, err)
}
