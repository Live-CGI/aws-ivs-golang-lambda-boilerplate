package main

import (
	dbConnection "aws-ivs-golang-serverless/db-connection"
	models "aws-ivs-golang-serverless/db-models"
	"aws-ivs-golang-serverless/utils"
	"context"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/uptrace/bun"
)

func getSearchParams(event events.APIGatewayV2HTTPRequest) (page int, pageSize int) {
	pageSizeStr := event.QueryStringParameters["pageSize"]
	pageStr := event.QueryStringParameters["page"]

	pageSize = 25
	page = 1

	if (pageSizeStr != "") {
		pageSizeInt, pageSizeErr := strconv.Atoi(pageSizeStr)

		if (pageSizeErr == nil) {
			pageSize = pageSizeInt
		}
	}

	if (pageStr != "") {
		pageInt, pageErr := strconv.Atoi(pageStr)

		if (pageErr == nil) {
			page = pageInt
		}
	}

	return page, pageSize
}

func getVideos(ctx context.Context, db *bun.DB, offset int, limit int) ([]models.OnDemandVideos, error) {
	results := make([]models.OnDemandVideos, 0)

	err := db.NewSelect().
		Model((*models.OnDemandVideos)(nil)).
		Where("active").
		Relation("UserContent").
		Limit(limit).
		Offset(offset).
		Scan(ctx, &results)

	return results, err
}

func buildResponsePayload(results []models.OnDemandVideos) ([]map[string]interface{}) {
	response := make([]map[string]interface{}, 0)

	// map structs to cleaned up json response
	for _, row := range results {
		// this is a generic interface - need to assert value types
		durationMs := row.Data["recording_duration_ms"].(float64)
		bucketName := row.Data["recording_s3_bucket_name"].(string)
		keyPrefix := row.Data["recording_s3_key_prefix"].(string)
		duration := durationMs / 1000

		content := map[string]interface{}{}

		if (row.UserContent != nil) {
			content = map[string]interface{}{
				"title": row.UserContent.Title,
				"description": row.UserContent.Description,
				"tags": row.UserContent.Tags,
				"thumbnail": row.UserContent.Thumbnail,
			}
		}

		response = append(response, map[string]interface{}{
			"uuid": row.Uuid,
			// created_at is when the recording ENDED, so start time is created_at row minus the duration miliseconds 
			"start_time": row.CreatedAt.UnixMilli() - int64(durationMs),
			"duration_seconds": duration,
			// TODO: after connecting cloudfront, replace bucketName + domain with CF domain
			"location": "https://" + bucketName + ".s3.amazonaws.com/" + keyPrefix + "/media/hls/master.m3u8",
			"active": row.Active,
			"content": content,
		})
	}

	return response
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (utils.Response, error) {
	// get request params
	page, pageSize := getSearchParams(event)

	db := dbConnection.GetDb()

	results, err := getVideos(ctx, db, (page - 1) * pageSize, pageSize)

	if (err != nil) {
		log.Println("Unable to get videos", err.Error())
		return utils.CreateErrorResponse(400, map[string]interface{}{ "error": err.Error() })
	}

	response := buildResponsePayload(results)

	res, err := utils.CreateOkResponse(200, map[string]interface{}{ "data": response, "page": page })

	return res, err
}

func main() {
	lambda.Start(Handler)
}
