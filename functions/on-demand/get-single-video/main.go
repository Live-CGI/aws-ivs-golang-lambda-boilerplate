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

func getOnDemandVideo(ctx context.Context, db *bun.DB, onDemandUuid uuid.UUID) (*models.OnDemandVideos, error) {
	onDemandVideo := new(models.OnDemandVideos)

	err := db.NewSelect().
		Model(onDemandVideo).
		Where("on_demand_videos.uuid = ?", onDemandUuid.String()).
		Relation("UserContent").
		Limit(1).
		Scan(ctx)
	
	if err != nil {
		log.Println("Error fetching video", err.Error())
		return nil, err
	}

	return onDemandVideo, nil
}

func parseRecordingJsonb(rec *models.OnDemandVideos) (float64, float64, string, string) {
	// ensure that the expected keys exist
	if _, ok := rec.Data["recording_duration_ms"]; ok {
		durationMs := rec.Data["recording_duration_ms"].(float64)
		bucketName := rec.Data["recording_s3_bucket_name"].(string)
		keyPrefix := rec.Data["recording_s3_key_prefix"].(string)
		durationS := durationMs / 1000
		return durationMs, durationS, bucketName, keyPrefix
	} else {
		// otherwise return with default values
		return 0, 0, "", ""
	}
}

func buildResponsePayload(rec *models.OnDemandVideos) (map[string]interface{}) {
	durationMs, duration, bucketName, keyPrefix := parseRecordingJsonb(rec)

	content := map[string]interface{}{}

	if (rec.UserContent != nil) {
		content = map[string]interface{}{
			"title": rec.UserContent.Title,
			"description": rec.UserContent.Description,
			"tags": rec.UserContent.Tags,
			"thumbnail": rec.UserContent.Thumbnail,
		}
	}

	response := map[string]interface{}{
		"uuid": rec.Uuid,
		"start_time": rec.CreatedAt.UnixMilli() - int64(durationMs),
		"duration_seconds": duration,
		// TODO: after connecting cloudfront, replace bucketName + domain with CF domain
		"location": "https://" + bucketName + ".s3.amazonaws.com/" + keyPrefix + "/media/hls/master.m3u8",
		"active": rec.Active,
		"content": content,
	}

	return response
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (utils.Response, error) {
	onDemandUuidStr := event.PathParameters["onDemandVideoUuid"]
	db := dbConnection.GetDb()

	onDemandUuid, err := uuid.Parse(onDemandUuidStr)

	if (err != nil) {
		log.Printf("Unparsable path params: %s", err.Error())
		return utils.CreateErrorResponse(400, map[string]interface{}{ "error": "Bad Request", "message": err.Error() })
	}

	video, err := getOnDemandVideo(ctx, db, onDemandUuid)

	if (err != nil) {
		log.Printf("Cannot find video: %s", err.Error())
		return utils.CreateErrorResponse(404, map[string]interface{}{ "error": "Recording not found" })
	}

	response := buildResponsePayload(video)

	return utils.CreateOkResponse(200, response);
}

func main() {
	lambda.Start(Handler)
}
