package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
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

func (p *SNSPublisher) PublishMessage(message string) error {
	_, err := p.snsClient.Publish(&sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(p.topicArn),
	})
	return err
}

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

func (r *SQSReceiver) ReceiveMessages() ([]string, error) {
	result, err := r.sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &r.queueUrl,
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(10),
	})

	if err != nil {
		return nil, err
	}

	messages := []string{}
	for _, msg := range result.Messages {
		messages = append(messages, *msg.Body)
	}
	return messages, nil
}
