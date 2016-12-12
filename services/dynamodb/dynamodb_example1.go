package main

// Created by Matt Harding
// Copyright go4aws.com

// Relied extensively on examples from:
// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#DynamoDB

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var tableName = "TableName"
var hashKeyName = "MyKey"
var service *dynamodb.DynamoDB

func main() {
	service = setupService()

	if !tableExists() {
		createTable()
	}
	addEntry()
	scanAndPrint()
	deleteTable()
}

func setupService() *dynamodb.DynamoDB {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return nil
	}

	svc := dynamodb.New(sess, aws.NewConfig().WithRegion("us-west-2"))
	return svc
}

func check(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
		return true
	} else {
		return false
	}
}

func tableExists() bool {

	params := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}
	resp, err := service.DescribeTable(params)

	if check(err) {
		return false
	}

	fmt.Println(resp)
	if *resp.Table.TableStatus == "ACTIVE" {
		return true
	}

	return false
}

func createTable() {
	fmt.Println("Creating...")

	params := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(hashKeyName),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(hashKeyName),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(tableName),
	}
	resp, err := service.CreateTable(params)

	if check(err) {
		return
	}

	fmt.Println(resp)

	// Wait until the table exists

	describeParams := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}

	err = service.WaitUntilTableExists(describeParams)
	check(err)
}

func deleteTable() {
	fmt.Println("Deleting...")

	input := &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	}
	resp, err := service.DeleteTable(input)
	if check(err) {
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)

	// Wait until the table no longer exists

	describeParams := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}

	err = service.WaitUntilTableNotExists(describeParams)
	if check(err) {
		return
	}
	fmt.Println("Done!")
}

func addEntry() {
	fmt.Println("Adding entry...")

	params := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			hashKeyName: {
				S: aws.String("MyValue"),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := service.PutItem(params)

	if check(err) {
		return
	}

	fmt.Println("Added entry.")
}

func scanAndPrint() {
	fmt.Println("Scanning...")

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	resp, err := service.Scan(input)

	if check(err) {
		return
	}

	fmt.Println(resp)
}
