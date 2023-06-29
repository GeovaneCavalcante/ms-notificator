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
	}

	sess := session.Must(session.NewSession(&cfg))

	svc := sns.New(sess)

	return &Sns{
		SnsConnection: svc,
		SnsArn:        snsArn,
	}
}

func (s *Sns) PublishMessage(message map[string]interface{}) (*messenger.MessageResponse, error) {

	msgRaw, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to marshal message: %v", err)
	}

	msgOutput, err := s.SnsConnection.Publish(&sns.PublishInput{
		Message:  aws.String(string(msgRaw)),
		TopicArn: &s.SnsArn,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to publish message: %v", err)
	}

	response := messenger.MessageResponse{
		ID: *msgOutput.MessageId,
	}

	return &response, err
}
