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

func toggleOnDemandVideo(ctx context.Context, db *bun.DB, onDemandVideoUuid uuid.UUID, on bool) (bool) {
	onDemandVideo := models.OnDemandVideos{
		Active: on,
		Uuid: onDemandVideoUuid,
	}

	_, err := db.NewUpdate().
		Model(&onDemandVideo).
		WherePK().
		Column("active").
		Returning("*").
		Exec(ctx)

	if (err != nil) {
		log.Printf("error toggling video: %s", err.Error())
		return false
	}

	return true
}

func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (utils.Response, error) {
	onDemandVideoUuidStr := event.PathParameters["onDemandVideoUuid"]
	db := dbConnection.GetDb()
	deleted := event.RequestContext.HTTP.Method == "DELETE"
	
	log.Printf("Request headers: %s", event.RequestContext.HTTP.Method)

	onDemandVideoUuid, recUuidErr := uuid.Parse(onDemandVideoUuidStr)

	log.Printf("video uuid %+v", onDemandVideoUuid)

	if (recUuidErr != nil) {
		log.Printf("Unparsable path params: %s", recUuidErr.Error())
		return utils.CreateErrorResponse(400, map[string]interface{}{ "error": "Bad Request", "message": recUuidErr.Error() })
	}

	ok := toggleOnDemandVideo(ctx, db, onDemandVideoUuid, !deleted)

	log.Printf("toggled ok %t", ok)
	
	if (!ok) {
		log.Println("failed to toggle video")
		return utils.CreateErrorResponse(500, map[string]interface{}{ "error": "Internal Server Error" })
	}

	return utils.CreateTextResponse(200, "OK");
}

func main() {
	lambda.Start(Handler)
}
