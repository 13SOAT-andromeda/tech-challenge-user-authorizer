#!/bin/bash
set -e

ENDPOINT_URL="http://localhost:4566"
FUNCTION_NAME="tech-challenge-user-authorizer"
ROLE_NAME="lambda-local-role"
REGION="us-east-1"

echo "Building Lambda..."
./scripts/build.sh

echo "Creating IAM Role in LocalStack..."
aws --endpoint-url=$ENDPOINT_URL iam create-role \
    --role-name $ROLE_NAME \
    --assume-role-policy-document '{"Version": "2012-10-17","Statement": [{ "Effect": "Allow", "Principal": {"Service": "lambda.amazonaws.com"}, "Action": "sts:AssumeRole"}]}' \
    --region $REGION || true

echo "Deploying Lambda to LocalStack..."
aws --endpoint-url=$ENDPOINT_URL lambda create-function \
    --function-name $FUNCTION_NAME \
    --runtime provided.al2023 \
    --role arn:aws:iam::000000000000:role/$ROLE_NAME \
    --handler bootstrap \
    --zip-file fileb://function.zip \
    --region $REGION \
    --environment "Variables={JWT_SECRET=test_secret,JWT_ISSUER=test_issuer}" || \
aws --endpoint-url=$ENDPOINT_URL lambda update-function-code \
    --function-name $FUNCTION_NAME \
    --zip-file fileb://function.zip \
    --region $REGION

echo "Deployment complete: $FUNCTION_NAME"
