package main

import (
	dbConnection "aws-ivs-golang-serverless/db-connection"
	models "aws-ivs-golang-serverless/db-models"
	"aws-ivs-golang-serverless/utils"
	"context"
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

func getIvsChannelState(ctx context.Context, db *bun.DB, uuid uuid.UUID) (*models.IvsStateChanges, error) {
	channelState := new(models.IvsStateChanges)

	err := db.NewSelect().
		Model(channelState).
		Where("ivs_channel_uuid = ?", uuid.String()).
		OrderExpr("timestamp DESC").
		Limit(1).
		Scan(ctx)
	
	if err != nil {
		log.Println("Error fetching ivs channel uuid", err.Error())
		return nil, err
	}

	return channelState, nil
}

func buildResponsePayload(state *models.IvsStateChanges) (map[string]interface{}) {
	response := map[string]interface{}{
		"timestamp": state.Timestamp.UnixMilli(),
		"state": state.State,
	}

	return response
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (utils.Response, error) {
	ivsChannelUuidStr := event.PathParameters["ivsChannelUuid"]
	db := dbConnection.GetDb()

	ivsChannelUuid, err := uuid.Parse(ivsChannelUuidStr)

	if (err != nil) {
		log.Printf("Unparsable path params: %s", err.Error())
		return utils.CreateErrorResponse(400, map[string]interface{}{ "error": "Bad Request", "message": err.Error() })
	}

	channelState, err := getIvsChannelState(ctx, db, ivsChannelUuid)

	if (err != nil) {
		log.Printf("Cannot find channel: %s", err.Error())
		return utils.CreateErrorResponse(404, map[string]interface{}{ "error": "IVS Channel not found" })
	}

	response := buildResponsePayload(channelState)

	return utils.CreateOkResponse(200, response);
}

func main() {
	lambda.Start(Handler)
}
