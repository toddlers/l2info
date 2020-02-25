package client

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func GetSession(region string) *session.Session {
	var err error
	if region == "" {
		region = "eu-central-1"
	}
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		fmt.Println("Failed to create session", err)
		return nil
	}
	return sess
}

func GetClient(session *session.Session) *lambda.Lambda {
	if session == nil {
		session = GetSession("")
	}
	svc := lambda.New(session)
	return svc
}
