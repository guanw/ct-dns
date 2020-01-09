package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guanw/ct-dns/storage"
)

// NewFactory creates dynamodb Client
func NewFactory() (storage.Client, error) {
	s := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:8000"),
	}))
	db := dynamodb.New(s)
	return NewClient(db), nil
}
