#!/bin/bash

# This script is used to upload the lambda function to AWS

# Get the info of the lambda function
lambdaName="gam-runner"

# Zip the lambda package folder
cd lambda_package && zip -r ../lambda.zip . && cd .. 

# Upload the lambda function to AWS
aws lambda update-function-code --function-name $lambdaName --zip-file fileb://lambda.zip > /dev/null

# Clean up the zip file
rm lambda.zip