package main

import (
	"errors"
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"strings"


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

	book := request.QueryStringParameters["book"]
	// If no name is provided in the HTTP request body, throw an error
	if len(book) < 1 {
		return events.APIGatewayProxyResponse{Body: ErrBookNotProvided.Error()}, ErrBookNotProvided
	}

	url := fmt.Sprintf("http://www.gutenberg.org/files/%s/%s.txt", book, book)
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	parsed := strings.Split(string(body), "\n")
	parsed = parsed[28:len(parsed)-398]

	return events.APIGatewayProxyResponse{
		Body:       strings.Join(parsed, "\n"),
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(Handler)
}
