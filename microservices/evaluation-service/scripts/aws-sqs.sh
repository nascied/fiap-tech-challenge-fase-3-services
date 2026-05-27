#!/bin/bash

#aguarda alguns segundo antes de criar a fila
sleep 15

#Cria fila sqs
awslocal sqs create-queue --queue-name evaluation-events

#Cria tabela no dynamodb
awslocal dynamodb create-table --table-name ToggleMasterAnalytics --attribute-definitions AttributeName=event_id,AttributeType=S --key-schema AttributeName=event_id,KeyType=HASH --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1