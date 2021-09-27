# simple-aws-sqs-app
A simple application that puts messages to a AWS SQS and reads them

## Localstack

### Start AWS localstack

docker-compose up -d localstack

### Stop AWS localstack

docker-compose down

### Run the application

go run .\cmd\app\main.go .\cmd\app\application.go

## AWS CLI commands

List all queues
### --endpoint-url=http://localhost:4566 sqs list-queues

### List queue attributes

aws --endpoint-url=http://localhost:4566 sqs get-queue-attributes --queue-url http://localhost:4566/000000000000/priority-spike-queue --attribute-names All

### Purge all messages in queue

aws --endpoint-url=http://localhost:4566 sqs purge-queue --queue-url http://localhost:4566/000000000000/priority-spike-queue
