package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type IvsChannel struct {
	bun.BaseModel 	`bun:"table:ivs_channels"`

	Uuid 			uuid.UUID `bun:"uuid,pk,default:generate_uuid_v4(),type:uuid"`
	Owner 			uuid.UUID `bun:"owner,notnull,type:uuid"`
	CreatedAt 		time.Time `bun:"created_at,notnull,default:now(),type:timestamptz"`
	Active 			bool `bun:"active,notnull,default:true"`
	Arn 			string `bun:"arn,notnull"`
	RtmpAddress 	string `bun:"rtmp_address,notnull"`
	StreamKey 		string `bun:"stream_key,notnull"`
	ChannelData 	map[string]interface{} `bun:"type:jsonb"`
}

type OnDemandVideos struct {
	bun.BaseModel 		`bun:"table:on_demand_videos"`

	Uuid 				uuid.UUID `bun:"uuid,pk,default:generate_uuid_v4(),type:uuid"`
	IvsChannelUuid 		uuid.UUID `bun:"ivs_channel_uuid,type:uuid"`
	CreatedAt 			time.Time `bun:"created_at,nullzero,notnull,default:now(),type:timestamptz"`
	Active 				bool `bun:"active,notnull,default:true"`
	Data 				map[string]interface{} `bun:"data,type:jsonb"`
	UserContent			*OnDemandVideoContent `bun:"rel:belongs-to,join:uuid=on_demand_video_uuid"`
}

type OnDemandVideoContent struct {
	bun.BaseModel 		`bun:"table:on_demand_video_content"`

	Uuid 				uuid.UUID `bun:"uuid,pk,default:generate_uuid_v4(),type:uuid"`
	OnDemandVideoUuid	uuid.UUID `bun:"on_demand_video_uuid,notnull,unique,type:uuid"`
	CreatedAt			time.Time `bun:"created_at,nullzero,notnull,default:now(),type:timestamptz"`
	CreatedBy			uuid.UUID `bun:"created_by,notnull,type:uuid"`
	Title				string `bun:"title,type:text"`
	Description			string `bun:"description,type:text"`
	Tags				[]string `bun:"tags,array"`
	Thumbnail			string `bun:"thumbnail"`
	Data				map[string]interface{} `bun:"data,type:jsonb"`
}

type IvsStateChanges struct {
	bun.BaseModel 		`bun:"table:ivs_state_changes"`

	Timestamp			time.Time `bun:"timestamp,nullzero,notnull,default:now(),type:timestamptz"`
	IvsChannelUuid 		uuid.UUID `bun:"ivs_channel_uuid,notnull,type:uuid"`
	State 				string `bun:"state,notnull"`
	Data 				map[string]interface{} `bun:"data,type:jsonb"`
}
