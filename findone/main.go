package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// List of movies
type Movie struct {
	ID    string
	Title string
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Connect to dynamoDB
	var sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Scan the table for all the values
	svc := dynamodb.New(sess)

	res, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(req.PathParameters["id"]),
			},
		},
	})

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       err.Error(),
		}, err
	}

	body, err := json.Marshal(Movie{
		ID:    *res.Item["ID"].S,
		Title: *res.Item["Title"].S,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       err.Error(),
		}, err
	}

	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(body),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
