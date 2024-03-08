package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	token := request.Headers["authorization"]
	tokenSlice := strings.Split(token, " ")
	authKey := os.Getenv("AUTH_API_KEY")
	
	var bearerToken string
	if len(tokenSlice) > 1 {
		bearerToken = tokenSlice[len(tokenSlice)-1]
	}

	if bearerToken != authKey {
		log.Println("Rejecting")
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{ IsAuthorized: false }, errors.New("Unauthorized")
	}

	log.Println("User is authorized")
	return events.APIGatewayV2CustomAuthorizerSimpleResponse{ IsAuthorized: true }, nil
}

func main() {
	lambda.Start(handler)
}
