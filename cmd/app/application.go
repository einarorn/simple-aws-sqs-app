package main

import (
	"fmt"
	"log"
	"simple-aws-sqs-app/internal/aws"
	"simple-aws-sqs-app/internal/aws/sqs"
	"time"
)

type App struct {
	QueueClient sqs.QueueClient
}

func NewApplication(addDelay bool) (App, error) {
	awsSession, err := aws.CreateAwsSession(aws.Config{
		Address: "http://localhost:4566",
		Region:  "eu-west-1",
		Profile: "localstack",
		ID:      "test",
		Secret:  "test",
	})

	if err != nil {
		return App{}, err
	}

	queueClient, err := sqs.NewQueueClient(sqs.NewQueueClientOptions{
		AwsSession:               awsSession,
		QueueName:                "priority-spike-queue",
		ReceiveWaitTimeSeconds:   10,
		VisibilityTimeoutSeconds: 30,
		MessageDelaySeconds:      0,
		ReceiveBatchSize:         10,
		AddDelay:                 addDelay,
	})

	app := App{QueueClient: queueClient}

	return app, err
}

func (a App) Prioritise1000Payments() error {
	var (
		firstCZKMessage time.Time
		lastCZKMessage  time.Time
		firstHUFMessage time.Time
		lastHUFMessage  time.Time
		firstEURMessage time.Time
		lastEURMessage  time.Time
	)

	paymentsCZK := 200
	paymentsHUF := 200
	paymentsEUR := 600

	fmt.Println("Adding CZK payments to queue...")
	for i := 0; i < paymentsCZK; i++ {
		err := a.QueueClient.SendMessage(fmt.Sprintf("CZK payment %d", i+1))

		if err != nil {
			return err
		}
	}

	fmt.Println("Adding HUF payments to queue...")
	for i := 0; i < paymentsHUF; i++ {
		err := a.QueueClient.SendMessage(fmt.Sprintf("HUF payment %d", i+1))

		if err != nil {
			return err
		}
	}

	fmt.Println("Adding EUR payments to queue...")
	for i := 0; i < paymentsEUR; i++ {
		err := a.QueueClient.SendMessage(fmt.Sprintf("EUR payment %d", i+1))

		if err != nil {
			return err
		}
	}

	counter := 0

	for {
		if !(counter < (paymentsCZK + paymentsHUF + paymentsEUR)) {
			break
		}

		batch, err := a.QueueClient.GetMessages()

		if err != nil {
			return err
		}

		for _, message := range batch.Messages {
			body := *message.Body
			log.Println(body)

			switch body[0:3] {
			case "CZK":
				if firstCZKMessage.IsZero() {
					firstCZKMessage = time.Now()
				}
				lastCZKMessage = time.Now()
			case "HUF":
				if firstHUFMessage.IsZero() {
					firstHUFMessage = time.Now()
				}
				lastHUFMessage = time.Now()
			case "EUR":
				if firstEURMessage.IsZero() {
					firstEURMessage = time.Now()
				}
				lastEURMessage = time.Now()
			}

			err := a.QueueClient.DeleteMessage(*message.ReceiptHandle)
			counter++

			if err != nil {
				return err
			}
		}
	}

	fmt.Println("\nFirst CZK payment received: ", firstCZKMessage.Format("15:04:05"))
	fmt.Println("Last CZK payment received:  ", lastCZKMessage.Format("15:04:05"))
	fmt.Println("First HUF payment received: ", firstHUFMessage.Format("15:04:05"))
	fmt.Println("Last HUF payment received:  ", lastHUFMessage.Format("15:04:05"))
	fmt.Println("First EUR payment received: ", firstEURMessage.Format("15:04:05"))
	fmt.Println("Last EUR payment received:  ", lastEURMessage.Format("15:04:05"))

	return nil
}
