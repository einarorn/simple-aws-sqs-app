package sqs

import "github.com/aws/aws-sdk-go/service/sqs"

type Queue interface {
	SendMessage(string) error
	GetMessages() (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(messageHandle string) error
}