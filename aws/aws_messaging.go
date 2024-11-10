package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/bernardbaker/qiba.core/ports"
)

type SNSPublisher struct {
	snsClient *sns.SNS
	topicArn  string
}

func NewSNSPublisher(topicArn string) *SNSPublisher {
	sess := session.Must(session.NewSession())
	return &SNSPublisher{
		snsClient: sns.New(sess),
		topicArn:  topicArn,
	}
}

func (p *SNSPublisher) PublishMessage(message ports.Message) error {
	_, err := p.snsClient.Publish(&sns.PublishInput{
		Message:  aws.String(message.Content),
		TopicArn: aws.String(p.topicArn),
	})
	return err
}

// Publish sends a message to the given viewers via SNS
func (p *SNSPublisher) Publish(content string, viewers []string) error {
	// Here, we'll send the content to SNS for the provided viewers.
	// For demonstration, we're just logging the message.
	for _, viewer := range viewers {
		// Example of publishing to SNS
		input := &sns.PublishInput{
			Message:  aws.String(fmt.Sprintf("Message for viewer %s: %s", viewer, content)),
			TopicArn: aws.String(p.topicArn),
		}

		_, err := p.snsClient.Publish(input)
		if err != nil {
			return fmt.Errorf("failed to publish message to SNS: %w", err)
		}
	}

	return nil
}

// SQSReceiver is responsible for receiving messages from SQS
type SQSReceiver struct {
	sqsClient *sqs.SQS
	queueUrl  string
}

func NewSQSReceiver(queueUrl string) *SQSReceiver {
	sess := session.Must(session.NewSession())
	return &SQSReceiver{
		sqsClient: sqs.New(sess),
		queueUrl:  queueUrl,
	}
}

// Receive receives messages from the SQS queue
func (r *SQSReceiver) Receive() ([]string, error) {
	// Receive messages from the SQS queue
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(r.queueUrl),
		MaxNumberOfMessages: aws.Int64(10), // Maximum number of messages to return
	}

	result, err := r.sqsClient.ReceiveMessage(input)
	if err != nil {
		return nil, fmt.Errorf("failed to receive message from SQS: %w", err)
	}

	messages := []string{}
	for _, message := range result.Messages {
		messages = append(messages, *message.Body) // Extracting message body
	}

	return messages, nil
}
