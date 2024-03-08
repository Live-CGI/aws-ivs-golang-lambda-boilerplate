package main

import (
	"aws-ivs-golang-serverless/utils"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ivs"
	iLambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/google/uuid"
)

type CreateChannelRequest struct {
	Name string `json:"name"`
	UserId string `json:"userId"`
}

type CreateChannelResponse struct {
	Arn string `json:"arn"`
	RtmpAddress string `json:"rtmpAddress"`
	StreamKey string `json:"streamKey"`
	Uuid string `json:"uuid"`
}

func createChannel(ctx context.Context, request *CreateChannelRequest, client *ivs.Client) (*ivs.CreateChannelOutput, error) {
	arn := os.Getenv("RECORDING_CONFIG_ARN")
	uuid := uuid.New().String()

	createChannelRequest := ivs.CreateChannelInput{
		Name: &uuid,
		RecordingConfigurationArn: &arn,
		Tags: map[string]string{
			// TODO remove illegal characters
			"Name": request.Name,
			"Owner": request.UserId,
		},
	}

	channel, err := client.CreateChannel(ctx, &createChannelRequest)

	if err != nil {
		log.Panic("error creating channel", err.Error())
		return nil, err
	}

	log.Printf("channel: %+v", channel)

	return channel, nil
}

func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (utils.Response, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }
	
	client := ivs.NewFromConfig(cfg)

	// get request json, throw if the json is not parseable
	request, reqErr := utils.ParseRequest[CreateChannelRequest](event)
	if reqErr != nil {
		log.Println("Error processing request", reqErr)
		return utils.CreateErrorResponse(500, map[string]interface{}{"error": reqErr.Error()})
	}

	// create the IVS channel
	channel, chErr := createChannel(ctx, &request, client)
	if chErr != nil {
		log.Println("Error creating channel", chErr)
		return utils.CreateErrorResponse(500, map[string]interface{}{"error": chErr.Error()})
	}

	// we want to store the IVS response body and this can be variable
	// so we should convert the unknown type to JSON for storage in DB
	channelDataBytes, chDataErr := json.Marshal(channel)
	if chDataErr != nil {
		log.Println("Error parsing channel data", chDataErr)
		return utils.CreateErrorResponse(500, map[string]interface{}{"error": chDataErr.Error()})
	}
	var channelData map[string]interface{}
	json.Unmarshal(channelDataBytes, &channelData)

	// create payload for inserting into the database
	writePayload, err := utils.CreateJsonPayload(map[string]interface{}{
		"channelArn": channel.Channel.Arn,
		"userId": request.UserId,
		"rtmpAddress": channel.Channel.IngestEndpoint,
		"streamKey": channel.StreamKey.Value,
		"channelData": channelData,
		"uuid": channel.Channel.Name,
	})

	if err != nil {
		log.Println("Error creating write payload", err)
		return utils.CreateErrorResponse(500, map[string]interface{}{"error": err.Error()})
	}

	// initialize lambda client
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println("Couldn't load default configuration. Have you set up your AWS account?")
		log.Println(err)
		return utils.CreateErrorResponse(500, map[string]interface{}{"error": err.Error()})
	}

	lambdaClient := iLambda.NewFromConfig(sdkConfig)
	writeFunc := os.Getenv("WRITE_FUNC")

	// invoke lambda
	// TODO: generate uuid on this side, insert uuid on write side, call async, respond to user
	// without awaiting response 
	invokeConfig := iLambda.InvokeInput{
		FunctionName: &writeFunc,
		Payload: writePayload,
	}

	out, err := lambdaClient.Invoke(ctx, &invokeConfig)
	if err != nil {
		log.Println("Error writing payload", err)
		return utils.CreateErrorResponse(500, map[string]interface{}{"error": err.Error()})
	}
	
	// parse json response from DB insert back to ensure payload is expected shape
	Json, err := utils.ParseJson[CreateChannelResponse](string(out.Payload))
	if err != nil {
		log.Println("Error parsing payload", err)
		return utils.CreateErrorResponse(500, map[string]interface{}{"error": err.Error()})
	}

	responseBody := map[string]interface{}{
		"streamKey": Json.StreamKey,
		"rtmpAddress": Json.RtmpAddress,
		"arn": Json.Arn,
		"uuid": Json.Uuid,
	}

	res, err := utils.CreateOkResponse(200, responseBody)
	return res, err
}

func main() {
	lambda.Start(Handler)
}