module aws-ivs-golang-serverless/service

go 1.20

replace aws-ivs-golang-serverless/utils => ./functions/modules/utils

replace aws-ivs-golang-serverless/db-models => ./functions/modules/db-models

replace aws-ivs-golang-serverless/db-connection => ./functions/modules/db-connection

require (
	aws-ivs-golang-serverless/db-connection v0.0.0-00010101000000-000000000000
	aws-ivs-golang-serverless/db-models v0.0.0-00010101000000-000000000000
	aws-ivs-golang-serverless/utils v0.0.0-00010101000000-000000000000
	github.com/aws/aws-lambda-go v1.46.0
	github.com/aws/aws-sdk-go-v2/config v1.27.6
	github.com/aws/aws-sdk-go-v2/service/ivs v1.33.1
	github.com/aws/aws-sdk-go-v2/service/lambda v1.53.1
	github.com/google/uuid v1.6.0
	github.com/uptrace/bun v1.1.17
)

require (
	github.com/aws/aws-sdk-go-v2 v1.25.2 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.1 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.6 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.15.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.20.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.23.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.28.3 // indirect
	github.com/aws/smithy-go v1.20.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/uptrace/bun/dialect/pgdialect v1.1.17 // indirect
	github.com/uptrace/bun/driver/pgdriver v1.1.17 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	mellium.im/sasl v0.3.1 // indirect
)
