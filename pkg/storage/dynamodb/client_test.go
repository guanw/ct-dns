package dynamodb

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guanw/ct-dns/pkg/storage/dynamodb/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {
	tests := []struct {
		Input       *dynamodb.PutItemInput
		ReturnErr   error
		Description string
		Value       string
		Key         string
		ExpectError bool
	}{
		{
			Input: &dynamodb.PutItemInput{
				Item: map[string]*dynamodb.AttributeValue{
					"Service": {
						S: aws.String("valid-service"),
					},
					"Host": {
						S: aws.String("192.0.0.1"),
					},
				},
				TableName: aws.String("service-discovery"),
			},
			ReturnErr:   nil,
			Description: "correctly set serviceName&host",
			Value:       "192.0.0.1",
			Key:         "valid-service",
			ExpectError: false,
		},
		{
			Input: &dynamodb.PutItemInput{
				Item: map[string]*dynamodb.AttributeValue{
					"Service": {
						S: aws.String("error-service"),
					},
					"Host": {
						S: aws.String("error"),
					},
				},
				TableName: aws.String("service-discovery"),
			},
			ReturnErr:   errors.New("error service"),
			Description: "set serviceName&host returns error",
			Value:       "error",
			Key:         "error-service",
			ExpectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			mockClient := &mocks.DynamodbClient{}
			c := NewClient(mockClient)
			mockClient.On("PutItem", test.Input).Return(&dynamodb.PutItemOutput{}, test.ReturnErr)
			err := c.Create(test.Key, test.Value)
			if test.ExpectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_Get(t *testing.T) {
	tests := []struct {
		Input       *dynamodb.QueryInput
		ReturnVal   *dynamodb.QueryOutput
		ReturnErr   error
		Description string
		Value       string
		Key         string
		ExpectError bool
	}{
		{
			Input: &dynamodb.QueryInput{
				TableName:              aws.String("service-discovery"),
				KeyConditionExpression: aws.String("Service = :service"),
				ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":service": {
						S: aws.String("empty-service"),
					},
				},
			},
			ReturnVal: &dynamodb.QueryOutput{
				Items: []map[string]*dynamodb.AttributeValue{},
			},
			Key:         "empty-service",
			Value:       `null`,
			ReturnErr:   nil,
			Description: "service name with no host registered",
			ExpectError: false,
		},
		{
			Input: &dynamodb.QueryInput{
				TableName:              aws.String("service-discovery"),
				KeyConditionExpression: aws.String("Service = :service"),
				ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":service": {
						S: aws.String("valid-service"),
					},
				},
			},
			ReturnVal: &dynamodb.QueryOutput{
				Items: []map[string]*dynamodb.AttributeValue{
					0: map[string]*dynamodb.AttributeValue{
						"Service": &dynamodb.AttributeValue{
							S: aws.String("valid-service"),
						},
						"Host": &dynamodb.AttributeValue{
							S: aws.String("192.0.0.1"),
						},
					},
				},
			},
			Key:         "valid-service",
			Value:       `["192.0.0.1"]`,
			ReturnErr:   nil,
			Description: "service name with one host registered",
			ExpectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			mockClient := &mocks.DynamodbClient{}
			c := NewClient(mockClient)
			mockClient.On("Query", test.Input).Return(test.ReturnVal, test.ReturnErr)
			val, err := c.Get(test.Key)
			if test.ExpectError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.Value, val)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_Delete(t *testing.T) {
	tests := []struct {
		Input       *dynamodb.DeleteItemInput
		ReturnErr   error
		Description string
		Value       string
		Key         string
		ExpectError bool
	}{
		{
			Input: &dynamodb.DeleteItemInput{
				Key: map[string]*dynamodb.AttributeValue{
					"Service": {
						S: aws.String("valid-service"),
					},
					"Host": {
						S: aws.String("192.0.0.1"),
					},
				},
				TableName: aws.String("service-discovery"),
			},
			ReturnErr:   nil,
			Description: "correctly delete serviceName&host",
			Value:       "192.0.0.1",
			Key:         "valid-service",
			ExpectError: false,
		},
		{
			Input: &dynamodb.DeleteItemInput{
				Key: map[string]*dynamodb.AttributeValue{
					"Service": {
						S: aws.String("error-service"),
					},
					"Host": {
						S: aws.String("error"),
					},
				},
				TableName: aws.String("service-discovery"),
			},
			ReturnErr:   errors.New("error service"),
			Description: "set serviceName&host returns error",
			Value:       "error",
			Key:         "error-service",
			ExpectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			mockClient := &mocks.DynamodbClient{}
			c := NewClient(mockClient)
			mockClient.On("DeleteItem", test.Input).Return(&dynamodb.DeleteItemOutput{}, test.ReturnErr)
			err := c.Delete(test.Key, test.Value)
			if test.ExpectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
