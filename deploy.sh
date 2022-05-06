#!/bin/sh
sam build

# Deploys with AuthDelay 100
sam deploy --no-confirm-changeset --no-fail-on-empty-changeset --stack-name "example-compare-ws-http-api" --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM --s3-bucket example-compare-ws-http-api --parameter-overrides ParameterKey=AuthDelay,ParameterValue=100

# Deploys with AuthDelay 0
#sam deploy --no-confirm-changeset --no-fail-on-empty-changeset --stack-name "example-compare-ws-http-api" --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM --s3-bucket example-compare-ws-http-api
