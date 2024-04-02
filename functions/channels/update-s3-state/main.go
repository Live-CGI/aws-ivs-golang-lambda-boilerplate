package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type StreamStateChangeDetail struct {
	EventName string `json:"event_name"`
	ChannelName string `json:"channel_name"`
	StreamId string `json:"stream_id"`
}

func putStateInS3(stage string, state string) (bool, error) {
	svc := s3.New(session.New())
	input := &s3.PutObjectInput{
		Body:   strings.NewReader("{ \"state\": \"" + state + "\" }"),
		Bucket: aws.String(os.Getenv("RECORDING_BUCKET")),
		Key:    aws.String(stage + "-live.json"),
	}
	result, err := svc.PutObject(input)

	log.Println("Put in s3 result", result)

	return err == nil, err;
}

func Handler(ctx context.Context, event events.EventBridgeEvent) (bool, error) {
	body := &StreamStateChangeDetail{}
	err := json.Unmarshal(event.Detail, body)

	if err != nil {
		log.Println("Error parsing json body", err.Error())
		return false, err
	}

	var jsonb map[string]interface{}
	jsonErr := json.Unmarshal(event.Detail, &jsonb)

	if jsonErr != nil {
		log.Println("Error parsing json body", jsonErr.Error())
		return false, jsonErr
	}

	var state string

	if body.EventName == "Stream Start" {
		state = "live"
	} else {
		state = "offline"
	}

	inS3, err := putStateInS3(os.Getenv("STAGE"), state);

	return inS3, err
}

func main() {
	lambda.Start(Handler)
}
