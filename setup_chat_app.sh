#!/bin/bash

# Function to create a file and add contents
create_file_with_content() {
    local filepath=$1
    local content=$2
    echo "$content" > "$filepath"
    echo "Created $filepath"
}

# Base project directory
PROJECT_DIR="chat-app"

# Create project structure
echo "Creating project directory structure..."

mkdir -p $PROJECT_DIR/{proto,domain,grpc,aws,repository}

# 1. Create proto/chat.proto
PROTO_CONTENT='syntax = "proto3";

package chat;

// The chat service definition.
service ChatService {
  // Sends a message from a viewer to the content creator
  rpc SendMessage (MessageRequest) returns (MessageResponse);
  
  // Content creator sends a message back to a single viewer or all viewers
  rpc SendBroadcast (BroadcastRequest) returns (BroadcastResponse);
}

// Message request from a viewer
message MessageRequest {
  string viewer_id = 1;
  string content = 2;
}

// Message response after a message is sent
message MessageResponse {
  string message_id = 1;
  string status = 2;
}

// Broadcast request from the content creator
message BroadcastRequest {
  string message = 1;
  repeated string viewers = 2; // list of viewer ids, if empty, send to all
}

// Broadcast response after broadcasting
message BroadcastResponse {
  string status = 1;
}'

create_file_with_content "$PROJECT_DIR/proto/chat.proto" "$PROTO_CONTENT"

# 2. Create domain/domain.go
DOMAIN_CONTENT='package domain

import "time"

type Message struct {
    ID        string
    ViewerID  string
    Content   string
    Timestamp int64
}

type ChatService interface {
    SendMessage(viewerID string, content string) (*Message, error)
    BroadcastMessage(content string, viewerIDs []string) error
}

type MessageRepository interface {
    SaveMessage(msg *Message) error
    GetMessages(viewerID string) ([]*Message, error)
}

func NewMessage(id, viewerID, content string) *Message {
    return &Message{
        ID:        id,
        ViewerID:  viewerID,
        Content:   content,
        Timestamp: time.Now().Unix(),
    }
}'

create_file_with_content "$PROJECT_DIR/domain/domain.go" "$DOMAIN_CONTENT"

# 3. Create grpc/grpc_service.go
GRPC_CONTENT='package grpc

import (
    "context"
    "github.com/google/uuid"
    "github.com/chat-app/domain"
    pb "github.com/chat-app/proto" // Proto generated code
)

type GRPCChatService struct {
    chatService domain.ChatService
}

func NewGRPCChatService(svc domain.ChatService) *GRPCChatService {
    return &GRPCChatService{svc}
}

func (s *GRPCChatService) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
    message, err := s.chatService.SendMessage(req.ViewerId, req.Content)
    if err != nil {
        return nil, err
    }
    return &pb.MessageResponse{
        MessageId: message.ID,
        Status:    "sent",
    }, nil
}

func (s *GRPCChatService) SendBroadcast(ctx context.Context, req *pb.BroadcastRequest) (*pb.BroadcastResponse, error) {
    err := s.chatService.BroadcastMessage(req.Message, req.Viewers)
    if err != nil {
        return nil, err
    }
    return &pb.BroadcastResponse{Status: "broadcasted"}, nil
}'

create_file_with_content "$PROJECT_DIR/grpc/grpc_service.go" "$GRPC_CONTENT"

# 4. Create aws/aws_messaging.go
AWS_CONTENT='package aws

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sns"
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

import (
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sqs"
)

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
}'

create_file_with_content "$PROJECT_DIR/aws/aws_messaging.go" "$AWS_CONTENT"

# 5. Create repository/memory_repository.go
REPO_CONTENT='package repository

import (
    "github.com/chat-app/domain"
)

type MemoryMessageRepository struct {
    messages []*domain.Message
}

func NewMemoryMessageRepository() *MemoryMessageRepository {
    return &MemoryMessageRepository{}
}

func (r *MemoryMessageRepository) SaveMessage(msg *domain.Message) error {
    r.messages = append(r.messages, msg)
    return nil
}

func (r *MemoryMessageRepository) GetMessages(viewerID string) ([]*domain.Message, error) {
    var result []*domain.Message
    for _, msg := range r.messages {
        if msg.ViewerID == viewerID {
            result = append(result, msg)
        }
    }
    return result, nil
}'

create_file_with_content "$PROJECT_DIR/repository/memory_repository.go" "$REPO_CONTENT"

# 6. Create main.go
MAIN_CONTENT='package main

import (
    "log"
    "net"

    "github.com/chat-app/aws"
    "github.com/chat-app/domain"
    "github.com/chat-app/grpc"
    "github.com/chat-app/repository"

    "google.golang.org/grpc"
    pb "github.com/chat-app/proto"
)

func main() {
    // Create the repository
    repo := repository.NewMemoryMessageRepository()

    // Create SNS Publisher and SQS Receiver
    snsPublisher := aws.NewSNSPublisher("sns-topic-arn")
    sqsReceiver := aws.NewSQSReceiver("sqs-queue-url")

    // Create the domain service
    chatService := domain.NewChatService(repo, snsPublisher, sqsReceiver)

    // Set up the gRPC server
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    chatGrpcService := grpc.NewGRPCChatService(chatService)
    pb.RegisterChatServiceServer(grpcServer, chatGrpcService)

    log.Printf("Server listening on port 50051...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}'

create_file_with_content "$PROJECT_DIR/main.go" "$MAIN_CONTENT"

# 7. Create CloudFormation template
CLOUDFORMATION_CONTENT='AWSTemplateFormatVersion: "2010-09-09"
Resources:
  ChatSNSTopic:
    Type: "AWS::SNS::Topic"
    Properties:
      TopicName: "ChatMessagesTopic"
  
  ChatSQSQueue:
    Type: "AWS::SQS::Queue"
    Properties: 
      QueueName: "ChatMessagesQueue"

  ChatSNSSQSPolicy:
    Type: "AWS::SQS::QueuePolicy"
    Properties:
      Queues:
        - !Ref ChatSQSQueue
      PolicyDocument:
        Version: "2012-10-17"
        Statement:
          - Effect: "Allow"
            Principal: "*"
            Action: "SQS:SendMessage"
            Resource: !GetAtt ChatSQSQueue.Arn
            Condition:
              ArnEquals:
                aws:SourceArn: !Ref ChatSNSTopic

  ChatSNSSubscription:
    Type: "AWS::SNS::Subscription"
    Properties:
      Endpoint: !GetAtt ChatSQSQueue.Arn
      Protocol: "sqs"
      TopicArn: !Ref ChatSNSTopic

Outputs:
  SNSTopicARN:
    Description: "The ARN of the SNS Topic"
    Value: !Ref ChatSNSTopic

  SQSQueueURL:
    Description: "The URL of the SQS Queue"
    Value: !Ref ChatSQSQueue
'

create_file_with_content "$PROJECT_DIR/cloudformation.yml" "$CLOUDFORMATION_CONTENT"

echo "Project setup completed!"
