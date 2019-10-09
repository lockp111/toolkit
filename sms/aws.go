package sms

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// AWSSender ...
type AWSSender struct {
	svc *sns.SNS
}

// NewAWSSender ...
func NewAWSSender(id, secret, token, regions string) *AWSSender {
	sess := session.New(&aws.Config{
		Region: aws.String(regions),
		Credentials: credentials.NewStaticCredentials(
			id,
			secret,
			token,
		),
	})

	return &AWSSender{
		svc: sns.New(sess),
	}
}

// Send ...
func (s *AWSSender) Send(phone, content string) (messageID string, err error) {
	var resp *sns.PublishOutput
	resp, err = s.svc.Publish(&sns.PublishInput{
		Message:     aws.String(content),
		PhoneNumber: aws.String(phone),
	})
	if err != nil {
		return
	}

	if resp.MessageId != nil {
		messageID = *resp.MessageId
	}
	return
}
