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
	// convert json input string to Record<string, any>
	var jsonb map[string]interface{}
	jsonErr := json.Unmarshal(event.Detail, &jsonb)

	if jsonErr != nil {
		log.Println("Error parsing json body", jsonErr.Error())
		return false, jsonErr
	}
	
	db := dbConnection.GetDb()

	channelName := jsonb["channel_name"].(string)

	ivsChannelUuid, err := getIvsChannel(ctx, db, channelName)

	if err != nil {
		log.Println("Error getting ivs channel", err.Error())
		return false, err
	}

	onDemandVideo := models.OnDemandVideos{
		IvsChannelUuid: *ivsChannelUuid,
		Active: true,
		Data: jsonb,
	}

	res, err := db.NewInsert().Model(&onDemandVideo).Exec(ctx)

	return res != nil, err

}

func main() {
	lambda.Start(Handler)
}
