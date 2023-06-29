package sns

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/GeovaneCavalcante/ms-notificator/config"
	"github.com/GeovaneCavalcante/ms-notificator/internal/messenger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type Sns struct {
	SnsConnection *sns.SNS
	SnsArn        string
}

func New(snsArn string) *Sns {

	envs := config.GetEnvVars()
	environment := envs.Environment
	region := envs.AwsRegion

	cfg := aws.Config{
		Region: aws.String(region),
	}

	if environment == "local" {
		awsEndpoint := envs.AwsBaseUrl
		cfg.Endpoint = aws.String(awsEndpoint)
		cfg.WithCredentials(credentials.AnonymousCredentials)
		log.Printf("Local environment detected. AWS Endpoint: %s", awsEndpoint)
	}

	sess := session.Must(session.NewSession(&cfg))

	svc := sns.New(sess)

	log.Printf("[Broker SNS] Successfully created SNS service")

	return &Sns{
		SnsConnection: svc,
		SnsArn:        snsArn,
	}
}

func (s *Sns) PublishMessage(message string) (*messenger.MessageResponse, error) {

	msgRaw, err := json.Marshal(message)
	if err != nil {
		log.Printf("[Broker SNS] Failed to marshal message: %v", err)
		return nil, fmt.Errorf("failed to marshal message: %v", err)
	}

	msgOutput, err := s.SnsConnection.Publish(&sns.PublishInput{
		Message:  aws.String(string(msgRaw)),
		TopicArn: &s.SnsArn,
	})

	if err != nil {
		log.Printf("[Broker SNS] Failed to publish message: %v", err)
		return nil, fmt.Errorf("failed to publish message: %v", err)
	}

	log.Printf("[Broker SNS] Successfully published message, MessageID: %s", *msgOutput.MessageId)

	response := messenger.MessageResponse{
		ID: *msgOutput.MessageId,
	}

	return &response, err
}
