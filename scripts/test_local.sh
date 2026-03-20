#!/bin/bash
set -e

ENDPOINT_URL="http://localhost:4566"
FUNCTION_NAME="tech-challenge-user-authorizer"
REGION="us-east-1"

echo "Generating valid JWT token..."
TOKEN=$(/usr/local/go/bin/go run scripts/gen_token.go)

echo "Test Case 1: Valid Token"
cat <<EOF > event.json
{
    "headers": {
        "Authorization": "Bearer $TOKEN"
    },
    "path": "/test"
}
EOF
aws --endpoint-url=$ENDPOINT_URL lambda invoke \
    --function-name $FUNCTION_NAME \
    --payload fileb://event.json \
    --region $REGION \
    response.json
cat response.json
echo ""

echo "Test Case 2: Invalid Token"
cat <<EOF > event.json
{
    "headers": {
        "Authorization": "Bearer invalid_token"
    },
    "path": "/test"
}
EOF
aws --endpoint-url=$ENDPOINT_URL lambda invoke \
    --function-name $FUNCTION_NAME \
    --payload fileb://event.json \
    --region $REGION \
    response.json
cat response.json
echo ""

echo "Test Case 3: Missing Header"
cat <<EOF > event.json
{
    "headers": {},
    "path": "/test"
}
EOF
aws --endpoint-url=$ENDPOINT_URL lambda invoke \
    --function-name $FUNCTION_NAME \
    --payload fileb://event.json \
    --region $REGION \
    response.json
cat response.json
echo ""

rm event.json response.json
