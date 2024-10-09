#!/bin/bash

# Your variables
function_name="gam-get-clients"

# Main Logic
export AWS_PAGER=""

GOOS=linux GOARCH=arm64 go build -o bootstrap main.go
zip function.zip bootstrap

aws lambda update-function-code --function-name "$function_name" \
--zip-file fileb://function.zip

# Clean up
rm bootstrap function.zip