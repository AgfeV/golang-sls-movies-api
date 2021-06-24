package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Movie struct {
	ID    string
	Title string
}

var svc *dynamodb.DynamoDB

func init() {

	// Connect to dynamoDB
	var sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Scan the table for all the values
	svc = dynamodb.New(sess)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (events.APIGatewayProxyResponse, error) {

	req, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while scanning DynamoTable" + err.Error(),
		}, nil
	}

	// Pull out only the data we need
	movies := make([]Movie, 0)
	for _, movie := range req.Items {
		movies = append(movies, Movie{
			ID:    *movie["ID"].S,
			Title: *movie["Title"].S,
		})
	}

	// JSON string of movies to return
	response, err := json.Marshal(movies)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 404}, err
	}

	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(response),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "world-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
