package main

import (
	"errors"
	"net/http"
	"io/ioutil"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// ErrNameNotProvided is thrown when a book id is not provided
	ErrBookNotProvided = errors.New("no query parameter 'book'")
)

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	// If no name is provided in the HTTP request body, throw an error
	if len(request.QueryStringParameters["book"]) < 1 {
		return events.APIGatewayProxyResponse{}, ErrBookNotProvided
	}

	resp, err := http.Get("http://www.gutenberg.org/files/1661/1661.txt")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)


	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(Handler)
}
