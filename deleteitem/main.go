package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Movie struct {
	ID    string
	Title string
}

func main() {
	lambda.Start(deleteItem)
}

func deleteItem(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Create a connection to the dynamodb
	var sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(req.PathParameters["id"]),
			},
		},
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		log.Fatalf("Got error calling DeleteItem: %s", err)
	}
	// PutItem to the dynamodb

	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "Sucessfully Deleted Item",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	return resp, nil
}
