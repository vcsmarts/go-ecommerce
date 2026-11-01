#!/bin/sh
# Runs once LocalStack is ready. Community edition does not persist state,
# so recreate the resources the app expects on every start.
set -eu

awslocal s3 mb s3://goecommerce-uploads || true
awslocal sqs create-queue --queue-name order-events