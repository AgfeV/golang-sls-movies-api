package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Movie struct {
	ID    string
	Title string
}

func main() {
	lambda.Start(putItem)
}

func putItem(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Grab the json request from the request body
	var movie Movie
	json.Unmarshal([]byte(req.Body), &movie)
	// Create a connection to the dynamodb
	var sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)
	av, err := dynamodbattribute.MarshalMap(movie)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
	// PutItem to the dynamodb

	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "Sucessfully PutItem",
	}
	return resp, nil
}
