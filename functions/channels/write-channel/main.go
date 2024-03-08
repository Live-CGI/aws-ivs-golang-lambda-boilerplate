package main

import (
	dbConnection "aws-ivs-golang-serverless/db-connection"
	models "aws-ivs-golang-serverless/db-models"
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

type InsertChannelEvent struct {
	ChannelArn string `json:"channelArn"`
	UserId string `json:"userId"`
	RtmpAddress string `json:"rtmpAddress"`
	StreamKey string `json:"streamKey"`
	ChannelData map[string]interface{} `json:"channelData"`
	Uuid uuid.UUID `json:"uuid"`
}

func insertChannel(ctx context.Context, event *InsertChannelEvent) (*models.IvsChannel, error) {
	db := dbConnection.GetDb()

	_ivsChannel := models.IvsChannel{
		Uuid: event.Uuid,
		Arn: event.ChannelArn,
		Owner: uuid.MustParse(event.UserId),
		RtmpAddress: event.RtmpAddress,
		StreamKey: event.StreamKey,
		ChannelData: event.ChannelData,
	}

	ivsChannel := &_ivsChannel

	_, err := db.NewInsert().Model(ivsChannel).Returning("*").Exec(ctx)

	if (err != nil) {
		return nil, err
	}

	return ivsChannel, nil
}

func Handler(ctx context.Context, event *InsertChannelEvent) (map[string]interface{}, error) {
	// insert IVS channel data into database
	ivsChannel, err := insertChannel(ctx, event)
	if err != nil {
		log.Println("Error inserting channel", err)
		return nil, err
	}

	// respond to request
	response := map[string]interface{}{
		"streamKey": ivsChannel.StreamKey,
		"rtmpAddress": ivsChannel.RtmpAddress,
		"arn": ivsChannel.Arn,
		"uuid": ivsChannel.Uuid.String(),
	}

	log.Println("Responding with payload", response)
	return response, err
}

func main() {
	lambda.Start(Handler)
}