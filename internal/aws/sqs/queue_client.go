package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
)

type QueueClient struct {
	queueUrl                 *string
	queueName                string
	receiveWaitTimeSeconds   *int64
	visibilityTimeoutSeconds *int64
	messageDelaySeconds      *int64
	receiveBatchSize         *int64
	sqsSvc                   *sqs.SQS
	addDelay                 bool
}

type NewQueueClientOptions struct {
	AwsSession               *session.Session
	QueueName                string
	ReceiveWaitTimeSeconds   int64
	VisibilityTimeoutSeconds int64
	MessageDelaySeconds      int64
	ReceiveBatchSize         int64
	AddDelay                 bool
}

func NewQueueClient(options NewQueueClientOptions) (QueueClient, error) {
	sqsSvc := sqs.New(options.AwsSession)
	queueUrl, err := createSqsQueue(sqsSvc, options.QueueName)
	if err != nil {
		return QueueClient{}, errors.Wrap(err, "create and get SQS queue url failed")
	}
	sqsClient := QueueClient{
		queueUrl:                 aws.String(queueUrl),
		queueName:                options.QueueName,
		receiveWaitTimeSeconds:   aws.Int64(options.ReceiveWaitTimeSeconds),
		visibilityTimeoutSeconds: aws.Int64(options.VisibilityTimeoutSeconds),
		receiveBatchSize:         aws.Int64(options.ReceiveBatchSize),
		messageDelaySeconds:      aws.Int64(options.MessageDelaySeconds),
		sqsSvc:                   sqsSvc,
		addDelay:                 options.AddDelay,
	}

	return sqsClient, nil
}

func (sqsClient *QueueClient) SendMessage(messageBody string) error {
	delaySeconds := int64(0)

	if sqsClient.addDelay {
		switch messageBody[0:3] {
		case "CZK":
			delaySeconds = int64(0)
		case "HUF":
			delaySeconds = int64(30)
		case "EUR":
			delaySeconds = int64(60)
		}
	}

		_, err := sqsClient.sqsSvc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:     sqsClient.queueUrl,
		MessageBody:  aws.String(messageBody),
		DelaySeconds: &delaySeconds,
	})
	return err
}

func (sqsClient *QueueClient) GetMessages() (*sqs.ReceiveMessageOutput, error) {
	msgResult, err := sqsClient.sqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            sqsClient.queueUrl,
		MaxNumberOfMessages: sqsClient.receiveBatchSize,
		VisibilityTimeout:   sqsClient.visibilityTimeoutSeconds, // avoid multiple consumers picking up the same message
		WaitTimeSeconds:     sqsClient.receiveWaitTimeSeconds,   // long-polling
	})
	if err != nil {
		return nil, err
	}

	return msgResult, nil
}

func (sqsClient *QueueClient) DeleteMessage(messageHandle string) error {
	_, err := sqsClient.sqsSvc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      sqsClient.queueUrl,
		ReceiptHandle: aws.String(messageHandle),
	})
	return err
}

func createSqsQueue(sqsSvc *sqs.SQS, queueName string) (string, error) {
	ret := "1209600" // 14 days

	res, err := sqsSvc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
		Attributes: map[string]*string{
			"MessageRetentionPeriod":        aws.String(ret),
			"VisibilityTimeout":             aws.String("30"),
			"ReceiveMessageWaitTimeSeconds": aws.String("10"),
		},
	})

	return *res.QueueUrl, err
}
