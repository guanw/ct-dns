#!/bin/bash

usage="Usage: $(basename "$0") [-h] [-e endpoint] [-r region] -- program to inject dynamodb schema for ct-dns

where:
    -h  show this help text
    -e  set the dynamodb endpoint (default: http://localhost:8000)
    -r  set the dynamodb region (default: us-east-1)"

endpoint="http://localhost:8000"
region="us-east-1"

while getopts ":he:r:" option; do
    case $option in
        h )  echo "$usage"
            exit 0
            ;;
        e )  endpoint=$OPTARG
            ;;
        r )  region=$OPTARG
            ;;
        \? ) printf "illegal option: -%s\n" "$OPTARG"
            echo "$usage"
            exit 1
            ;;
    esac
done
shift $((OPTIND-1))
echo "Creating ct-dns table in dynamodb"
echo "set endpoint: $endpoint"
echo "set region: $region"


aws dynamodb --endpoint-url "$endpoint" --region "$region" \
	create-table \
	--table-name service-discovery \
    --attribute-definitions AttributeName=Service,AttributeType=S AttributeName=Host,AttributeType=S \
	--key-schema AttributeName=Service,KeyType=HASH AttributeName=Host,KeyType=RANGE \
	--provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1