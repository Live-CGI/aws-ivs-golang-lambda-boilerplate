package utils

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type Response events.APIGatewayProxyResponse

func CreateJsonPayload(Body map[string]interface{}) ([]byte, error) {
	body, err := json.Marshal(Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func CreateJsonStringPayload(Body map[string]interface{}) (*string, error) {
	body, err := json.Marshal(Body)
	if err != nil {
		return nil, err
	}
	str := string(body)
	return &str, nil
}

func CreateRawResponse(StatusCode int, Body []byte) (Response, error) {
	resp := Response{
		StatusCode: StatusCode,
		IsBase64Encoded: true,
		Body: string(Body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	return resp, nil
}

func CreateOkResponse(StatusCode int, Body map[string]interface{}) (Response, error) {
	var buf bytes.Buffer

	body, err := json.Marshal(Body)
	if err != nil {
		return Response{StatusCode: 500}, err
	}
	json.HTMLEscape(&buf, body)
	resp := Response{
		StatusCode: StatusCode,
		IsBase64Encoded: false,
		Body: buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	return resp, nil
}

func CreateTextResponse(StatusCode int, Body string) (Response, error) {
	resp := Response{
		StatusCode: StatusCode,
		IsBase64Encoded: false,
		Body: Body,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}
	return resp, nil
}

func CreateErrorResponse(StatusCode int, Body map[string]interface{}) (Response, error) {
	var buf bytes.Buffer

	body, err := json.Marshal(Body)
	if err != nil {
		return Response{StatusCode: 500}, err
	}
	json.HTMLEscape(&buf, body)
	resp := Response{
		StatusCode: StatusCode,
		IsBase64Encoded: false,
		Body: buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	return resp, nil
}

func ParseRequest[T any](event events.APIGatewayV2HTTPRequest) (T, error) {
	body := []byte(event.Body)
	var request T

	jsonErr := json.Unmarshal(body, &request)
	if (jsonErr != nil) {
		return request, jsonErr
	} else {
		return request, nil
	}
}

func ParseJson[T any](payload string) (*T, error) {
	bytes := []byte(payload)
	var Json *T
	err := json.Unmarshal(bytes, &Json)

	if (err != nil) {
		return nil, err
	}

	return Json, nil
}