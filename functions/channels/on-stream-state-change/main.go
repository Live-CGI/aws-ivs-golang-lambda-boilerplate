package main

import (
	dbConnection "aws-ivs-golang-serverless/db-connection"
	models "aws-ivs-golang-serverless/db-models"
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type StreamStateChangeDetail struct {
	EventName string `json:"event_name"`
	ChannelName string `json:"channel_name"`
	StreamId string `json:"stream_id"`
}

func getIvsChannel(ctx context.Context, db *bun.DB, uuidString string) (*uuid.UUID, error) {
	ivsChannelUuid := &uuid.UUID{}

	err := db.NewSelect().
		Model((*models.IvsChannel)(nil)).
		Column("uuid").
		Where("uuid = ?", uuidString).
		Limit(1).
		Scan(ctx, ivsChannelUuid)
	
	if err != nil {
		log.Println("Error fetching ivs channel uuid", err.Error())
		return nil, err
	}

	return ivsChannelUuid, nil
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
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

	db := dbConnection.GetDb()

	ivsChannelUuid, err := getIvsChannel(ctx, db, body.ChannelName)

	if err != nil {
		log.Println("Error getting ivs channel", err.Error())
		return false, err
	}

	var state string

	if body.EventName == "Stream Start" {
		state = "live"
	} else {
		state = "offline"
	}

	ivsStreamStateChange := models.IvsStateChanges{
		IvsChannelUuid: *ivsChannelUuid,
		State: state,
		Data: jsonb,
	}

	dbInsert, err := db.NewInsert().Model(&ivsStreamStateChange).Exec(ctx)

	log.Println(dbInsert)

	return err != nil, err
}

func main() {
	lambda.Start(Handler)
}
