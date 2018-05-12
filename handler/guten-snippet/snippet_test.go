package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)


type ClientMock struct {
}

func (c *ClientMock) Get(url string) (*http.Response, error) {
	testData, _ := ioutil.ReadFile("1661.txt")
    return &http.Response{Body: ioutil.NopCloser(bytes.NewReader(testData))}, nil
}

func TestHandler(t *testing.T) {
	mockClient := &ClientMock{}

	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  string
		err     error
	}{
		{
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{"book":"1661","start":"6","limit":"1"}},
			expect:  "THE ADVENTURES OF SHERLOCK HOLMES",
			err:     nil,
		},
		{
			request: events.APIGatewayProxyRequest{Body: ""},
			expect:  ErrBookNotProvided.Error(),
			err:     ErrBookNotProvided,
		},
	}

	for _, test := range tests {
		response, err := HandleRequest(test.request, mockClient)
		log.Printf(response.Body)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}

}
