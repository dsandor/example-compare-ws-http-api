#!/bin/sh
sam build
sam deploy --no-confirm-changeset --no-fail-on-empty-changeset --stack-name "example-compare-ws-http-api" --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM --s3-bucket example-compare-ws-http-api
