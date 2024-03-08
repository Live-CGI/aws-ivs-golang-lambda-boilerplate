--CREATE DATABASE ivs;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS ivs_channels (
    "uuid" uuid DEFAULT uuid_generate_v4(),
    "owner" uuid NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT now(),
    "active" BOOLEAN DEFAULT true,
    "arn" TEXT NOT NULL,
    "rtmp_address" TEXT NOT NULL,
    "stream_key" TEXT NOT NULL,
    "channel_data" JSONB NOT NULL,
    PRIMARY KEY ("uuid")
);

CREATE INDEX IF NOT EXISTS 
    idx_ivs_channels_owner ON 
    ivs_channels 
    USING HASH ("owner");

CREATE TABLE IF NOT EXISTS ivs_state_changes (
    "timestamp" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "ivs_channel_uuid" uuid NOT NULL,
    "state" TEXT NOT NULL,
    "data" JSONB,
    CONSTRAINT fk_ivs_state_changes_to_ivs_channel 
        FOREIGN KEY ("ivs_channel_uuid")
        REFERENCES ivs_channels("uuid")
);

CREATE INDEX IF NOT EXISTS 
    idx_ivs_state_channels_channel ON 
    ivs_state_changes
    USING BTREE ("ivs_channel_uuid", "timestamp");

CREATE TABLE IF NOT EXISTS on_demand_videos (
    "uuid" uuid DEFAULT uuid_generate_v4(),
    "ivs_channel_uuid" uuid,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "active" BOOLEAN DEFAULT true,
    "data" JSONB NOT NULL,
    PRIMARY KEY ("uuid")
);

CREATE TABLE IF NOT EXISTS on_demand_video_content (
    "uuid" uuid DEFAULT uuid_generate_v4(),
    "on_demand_video_uuid" uuid NOT NULL UNIQUE,
    "created_at" TIMESTAMPTZ DEFAULT now(),
    "created_by" uuid NOT NULL,
    "title" TEXT,
    "description" TEXT,
    "tags" _TEXT,
    "data" JSONB,
    "thumbnail" TEXT,
    PRIMARY KEY ("uuid"),
    CONSTRAINT fk_user_content_for_on_demand_video
        FOREIGN KEY ("on_demand_video_uuid")
        REFERENCES on_demand_videos("uuid")
);

CREATE INDEX IF NOT EXISTS 
    idx_on_demand_video_user_content ON 
    on_demand_video_content 
    USING HASH ("on_demand_video_uuid");