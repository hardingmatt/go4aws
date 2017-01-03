package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var service *sns.SNS
var topicName = "myTopic"

func main() {
	service = setupService()
	if !topicExists(topicName) {
		createTopic(topicName)
	}

	publishTextMessage("+18608087099", "Test")
	publishToTopic(topicName, "Test")

	deleteTopic(topicName)
}

func setupService() *sns.SNS {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return nil
	}

	svc := sns.New(sess, aws.NewConfig().WithRegion("us-west-2"))
	return svc
}

func topicExists(topicName string) bool {
	input := &sns.GetTopicAttributesInput{
		TopicArn: aws.String(createTopicArn(topicName)),
	}
	_, err := service.GetTopicAttributes(input)
	return err == nil
}

func createTopic(topicName string) {
	input := &sns.CreateTopicInput{
		Name: aws.String(topicName),
	}
	service.CreateTopic(input)
}

func deleteTopic(topicName string) {
	input := &sns.DeleteTopicInput{
		TopicArn: aws.String(createTopicArn(topicName)),
	}
	service.DeleteTopic(input)
}

func publishTextMessage(phoneNumber, message string) {
	input := &sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String(phoneNumber),
	}
	resp, err := service.Publish(input)
	fmt.Println(resp, err)
}

func publishToTopic(topicName, message string) {
	input := &sns.PublishInput{
		TopicArn: aws.String(createTopicArn(topicName)),
		Message:  aws.String(message),
	}
	resp, err := service.Publish(input)
	fmt.Println(resp, err)
}

func createTopicArn(topicName string) string {
	return "arn:aws:sns:us-west-2:311139150838:" + topicName
}
