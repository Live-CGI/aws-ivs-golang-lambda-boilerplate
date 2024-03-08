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

func getSearchParams(event events.APIGatewayV2HTTPRequest) (page int, pageSize int, search string) {
	pageSizeStr := event.QueryStringParameters["pageSize"]
	pageStr := event.QueryStringParameters["page"]
	search = event.QueryStringParameters["search"]

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

	return page, pageSize, search
}

func searchOnDemandContent(ctx context.Context, db *bun.DB, search string, offset int, limit int) ([]models.OnDemandVideoContent, error) {
	content := make([]models.OnDemandVideoContent, 0)

	err := db.NewRaw(`
		SELECT 
			on_demand_video_uuid,
			title,
			description,
			thumbnail,
			tags
			-- build from subquery so we can search up and down the result set
			-- while preserving search order 
			FROM (
				SELECT 
					*,
					-- build a search vector from title, description and tags
					to_tsvector(
						CONCAT(
							title, ' ', 
							description, ' ', 
							array_to_string(tags, ' ')
						)
					) search,
					websearch_to_tsquery(?) querytext
				FROM 
					on_demand_video_content uc
			) searchable
		WHERE querytext @@ search 
		ORDER BY ts_rank_cd(search, querytext) DESC 
		OFFSET ? LIMIT ?`, 
		search, offset, limit).
		Scan(ctx, &content)
	
	if err != nil {
		log.Println("Error fetching vidoes", err.Error())
		return nil, err
	}

	return content, nil
}

func buildResponsePayload(results []models.OnDemandVideoContent) ([]map[string]interface{}) {
	response := make([]map[string]interface{}, 0)

	for _, row := range results {
		response = append(response, map[string]interface{}{
			"title": row.Title,
			"description": row.Description,
			"tags": row.Tags,
			"thumbnail": row.Thumbnail,
			"on_demand_video_uuid": row.OnDemandVideoUuid,
		})
	}

	return response
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (utils.Response, error) {
	page, pageSize, search := getSearchParams(event)

	log.Printf("Searching for %s starting on page %d limit %d", search, page, pageSize)

	db := dbConnection.GetDb()

	videos, err := searchOnDemandContent(ctx, db, search, (page - 1) * pageSize, pageSize)

	if (err != nil) {
		log.Printf("Cannot find videos: %s", err.Error())
		return utils.CreateErrorResponse(404, map[string]interface{}{ "error": "Videos not found" })
	}

	response := buildResponsePayload(videos)

	res, err := utils.CreateOkResponse(200, map[string]interface{}{ "data": response, "page": page })

	return res, err
}

func main() {
	lambda.Start(Handler)
}
