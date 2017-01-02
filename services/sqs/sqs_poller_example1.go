package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var service *sqs.SQS

var queueName = "goSqs"

func main() {
	service = setupService()

	queueUrl := nameToUrl("us-west-2", "311139150838", queueName)
	if !queueExists(queueUrl) {
		createQueue(queueName)
	}
	poll(queueUrl)
}

func setupService() *sqs.SQS {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return nil
	}

	svc := sqs.New(sess, aws.NewConfig().WithRegion("us-west-2"))
	return svc
}

func nameToUrl(region, accountId, queueName string) string {
	return "https://sqs." + region + ".amazonaws.com/" + accountId + "/" + queueName
}

func queueExists(queueUrl string) bool {
	input := &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(queueUrl),
	}
	_, err := service.GetQueueAttributes(input)

	exists := err == nil
	if exists {
		fmt.Println("Queue exists: " + queueUrl)
	} else {
		fmt.Println("Queue does not exist: " + queueUrl)
	}
	return exists
}

func createQueue(queueName string) bool {
	input := &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	}
	resp, err := service.CreateQueue(input)

	success := err == nil
	if success {
		fmt.Println("Created queue: ", resp)
	} else {
		fmt.Println("Failed to create queue: ", err)
	}
	return success
}

func poll(queueUrl string) {
	fmt.Println("Polling...")

	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueUrl),
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(2),
	}
	resp, err := service.ReceiveMessage(input)
	fmt.Println("Poll results:", resp, "\nError:", err)
}
