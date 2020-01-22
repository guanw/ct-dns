package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guanw/ct-dns/pkg/logging"
	"github.com/guanw/ct-dns/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type builder struct {
	Region   string
	Endpoint string
}

func initFromViper(v *viper.Viper) *builder {
	return &builder{
		Region:   v.GetString("dynamodb-region"),
		Endpoint: v.GetString("dynamodb-endpoint"),
	}
}

// NewFactory creates dynamodb Client
func NewFactory(v *viper.Viper) (storage.Client, error) {
	b := initFromViper(v)
	s := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(b.Region),
		Endpoint: aws.String(b.Endpoint),
	}))
	logging.GetLogger().WithFields(logrus.Fields{
		"Endpoint": b.Endpoint,
		"Region":   b.Region,
	}).Info("Creating dynamodb session")
	db := dynamodb.New(s)
	return NewClient(db), nil
}
