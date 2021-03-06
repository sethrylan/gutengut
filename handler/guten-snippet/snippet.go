package main

import (
	"errors"
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

var (
	// ErrNameNotProvided is thrown when a book id is not provided
	ErrBookNotProvided = errors.New("no query parameter 'book'.")
	ErrOutOfBounds = errors.New("'start' and 'limit' out of bounds.")
)

type HttpClient interface {
    // Do(req *http.Request) (*http.Response, error)
    Get(url string) (*http.Response, error)
}

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing request %s\n", request.RequestContext.RequestID)	// stdout and stderr are sent to AWS CloudWatch Logs

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	return HandleRequest(request, httpClient)
}

func main() {
	lambda.Start(Handler)
}

func HandleRequest(request events.APIGatewayProxyRequest, httpClient HttpClient) (events.APIGatewayProxyResponse, error) {
	book := request.QueryStringParameters["book"]						// Unpack query params;
	start, _ := strconv.Atoi(request.QueryStringParameters["start"])	// default to 0 for start
	limit, _ := strconv.Atoi(request.QueryStringParameters["limit"])	// and limit

	if len(book) < 1 {													// If no name is provided in the HTTP request body, throw an error
		return events.APIGatewayProxyResponse{Body: ErrBookNotProvided.Error()}, ErrBookNotProvided
	}

	url := fmt.Sprintf("http://www.gutenberg.org/files/%s/%s.txt", book, book)
	resp, err := httpClient.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)			// Read response body
	parsed := strings.Split(string(body), "\n")	    // and converse to array of lines
	parsed = parsed[28:len(parsed)-398]				// and skip PG boilerplate lines (first 28 and last 398)

	if start + limit > len(parsed) {
		return events.APIGatewayProxyResponse{Body: ErrOutOfBounds.Error()}, ErrOutOfBounds
	}

	if limit == 0 {									// if limit is 0
		parsed = parsed[start:]						// then treat as no limit
	} else {										// else
		parsed = parsed[start:start+limit]			// use the passed limit
	}

	return events.APIGatewayProxyResponse {
		Body:       strings.Join(parsed, "\n"),
		StatusCode: 200,
		Headers:    map[string]string{"content-type": "text/plain"},
	}, nil
}
