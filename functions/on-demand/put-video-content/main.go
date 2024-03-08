package main

import (
	dbConnection "aws-ivs-golang-serverless/db-connection"
	models "aws-ivs-golang-serverless/db-models"
	"aws-ivs-golang-serverless/utils"
	"context"
	"database/sql"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PutVideoContentRequest struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Thumbnail string `json:"thumbnail"`
	Tags []string `json:"tags"`
}

func putVideoContent(ctx context.Context, db *bun.DB, onDemandVideoUuid uuid.UUID, request PutVideoContentRequest) (models.OnDemandVideoContent, sql.Result, error) {
	onDemandVideoContent := models.OnDemandVideoContent{
		OnDemandVideoUuid: onDemandVideoUuid,
		Title: request.Title,
		Description: request.Description,
		Thumbnail: request.Thumbnail,
		Tags: request.Tags,
	}

	rows, err := db.NewInsert().
		Model(&onDemandVideoContent).
		On("CONFLICT (on_demand_video_uuid) DO UPDATE").
		Set("title = EXCLUDED.title").
		Set("description = EXCLUDED.description"). 
		Set("thumbnail = EXCLUDED.thumbnail").
		Set("tags = EXCLUDED.tags").
		Returning("*").
		Exec(ctx)
	
	return onDemandVideoContent, rows, err
}

func buildResponsePayload(onDemandVideoContent models.OnDemandVideoContent) (map[string]interface{}) {
	response := map[string]interface{}{
		"uuid": onDemandVideoContent.Uuid.String(),
		"on_demand_video_uuid": onDemandVideoContent.OnDemandVideoUuid.String(),
		"title": onDemandVideoContent.Title,
		"description": onDemandVideoContent.Description,
		"thumbnail": onDemandVideoContent.Thumbnail,
		"tags": onDemandVideoContent.Tags,
		"created_at": onDemandVideoContent.CreatedAt,
	}

	return response
}

func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (utils.Response, error) {
	onDemandVideoUuidStr := event.PathParameters["onDemandVideoUuid"]
	db := dbConnection.GetDb()

	onDemandVideoUuid, recUuidErr := uuid.Parse(onDemandVideoUuidStr)

	if (recUuidErr != nil) {
		log.Printf("Unparsable path params: %s", recUuidErr.Error())
		return utils.CreateErrorResponse(400, map[string]interface{}{ "error": "Bad Request", "message": recUuidErr.Error() })
	}

	request, err := utils.ParseRequest[PutVideoContentRequest](event)

	if (err != nil) {
		log.Printf("Unparsable request body, error: %s", err.Error())
		return utils.CreateErrorResponse(400, map[string]interface{}{ "error": "Bad Request", "message": err.Error() })
	}
	
	onDemandVideoContent, rows, insErr := putVideoContent(ctx, db, onDemandVideoUuid, request)
	
	if (insErr != nil) {
		log.Println("failed to insert post", insErr)
		return utils.CreateErrorResponse(500, map[string]interface{}{ "error": "Internal Server Error" })
	}

	log.Println("Inserted", rows, "rows")

	response := buildResponsePayload(onDemandVideoContent)

	return utils.CreateOkResponse(200, response)
}

func main() {
	lambda.Start(Handler)
}
