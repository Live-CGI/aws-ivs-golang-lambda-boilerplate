# AWS IVS Golang Serverless Boilerplate

This repo is intended as a boilerplate for setting up an IVS channel or channels with recording and the ability to list, search and manage recorded videos. 

### ⚠️ Please note ⚠️

> This boilerplate is intended to be a **starting point**, not a fully production-ready service. The authorization method is a simple API key and should not be considered fully secure, and there is no security method at all on most resources (getting recordings, checking if a channel is live, etc.) There is also no security on accessing the live stream or the videos themselves, nor is any authentication method provided. The bucket created for recordings will not be generally accessible for video playback on the internet, and it is recommended that you attach a Cloudfront distro to the bucket if you want to use it for playback. 

## Setup

### Things you need before you start
1. An AWS Account
2. The AWS cli installed and configured for the account
3. A Postgres database in RDS available in at least one subnet in your VPC
4. One or more security groups needed to connect to the RDS database from AWS Lambda
5. NodeJS and Go installed
6. Serverless Framework installed `npm install -g serverless@3` 

### Things you need to do before you can deploy
1. Create a dotenv file at the project root named `.env` with the following keys:
```
# hostname of your postgresql database
PG_HOST=myHostname

# database name
PG_DATABASE=myDatabase

# database username
PG_USER=myUser

# database password
PG_PASSWORD=myPassword

# id of your subnet which can connect to postgres
SUBNET_1=subnet-0000000000000000

# id of your security group granting access to postgres
SECURITY_GROUP=sg-0000000000000001

# other ID - add as many as needed in list below
SECURITY_GROUP_2=sg-0000000000000002

# API key for protected routes
AUTH_API_KEY=my-super-secret-api-key
```

2. Run the SQL file at /db-bootstrap/init.sql to create the tables
    - uncomment the first line to create a new databse within Postgres
    - if you'd like to rename these tables, ensure that the configs in 
    /functions/modules/db-models is updated to match the new configuration

## Deploying
If you run `serverless deploy` it will deploy to your default AWS account in us-east-1 with the stage "dev", but common options are below:
```
serverless deploy \
    --aws-profile myOptionalProfile
    --region some-aws-region
    --stage someStage
```

## Rest Endpoints Created:
```
  POST      /channels
  GET       /on-demand
  GET       /on-demand/search
  PUT       /on-demand/{onDemandVideoUuid}
  GET       /on-demand/{onDemandVideoUuid}
  DELETE    /on-demand/{onDemandVideoUuid}
  PATCH     /on-demand/{onDemandVideoUuid}
```

## Scale considerations:
- 🚨 Pay attention to your resource costs 🚨 - IVS is not a cheap service at scale, nor is RDS. If you provide this out to the internet at large **Your bill will be expensive**. 
- No cache layer is provided. Each GET request for recordings or the live state will hit the database directly, and in a high-volume situation adding some cache in Redis, Memcached or even at the domain level will keep database usage low.
- I have included some basic indexes on the database tables, but these are based on some general rules based on how I've chosen to access data in these tables. As you extend functionality these indexes might need extending, modifying, or even potentially dropping.
- I have included some JSONB columns in the databases simply to hold onto all the data provided by IVS but if you want to depend on any of the columns in your sort or selection logic I would recommend adding new columns to the tables to contain that data separately. 

## Security considerations:
- The authorization method used is only intended to provide minimal functionality and should be improved to use something standards-based
- The AWS Resources and IAM roles in the serverless.yml file only provide minimal functionality and can/should be extended for your own security requirements. The roles created have very open rules.
- To be able to communicate both with the outside internet AND the VPC containing the RDS database, you would generally create a NAT Gateway as by default a VPC cannot connect to things outside of the VPC. NAT Gateways cost money and because this is used by me as a personal boilerplate not an organizational one I've put the IVS channel creation lambda (which needs access outside the VPC) and the lambda that writes the data to the database in separate lambdas: one inside the VPC and one outside the VPC. Your own intrastructure requirements might (and probably will) differ. This is just the cheap method, not necessarily the best method. 

Contributions, reviews and critiques are welcome!