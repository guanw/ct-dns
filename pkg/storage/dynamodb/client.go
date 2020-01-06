package dynamodb

import (
	"encoding/json"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/clicktherapeutics/ct-dns/pkg/storage"
	"github.com/pkg/errors"
)

// Client defines api client for set/get operations
type Client struct {
	DB   DynamodbClient
	lock sync.Mutex
}

// DynamodbClient defines the interface for dynamodb client
type DynamodbClient interface {
	Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	PutItem(putItemInput *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	DeleteItem(deleteItemInput *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
}

// Params defines config to initialize dynamodb client
type Params struct {
	Endpoint string
	Region   string
}

// NewClient creates new api client
func NewClient(db DynamodbClient) storage.Client {
	return &Client{
		DB: db,
	}
}

type keyValuePair struct {
	Service string `dynamodbav: "Service"`
	Host    string `dynamodbav: "Host"`
}

// CreateOrSet create new key/value pair or set existing key
func (c *Client) Create(key, value string) error {
	s := keyValuePair{
		Service: key,
		Host:    value,
	}
	sMap, err := dynamodbattribute.MarshalMap(s)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal serviceToHost map")
	}
	c.lock.Lock()
	_, err = c.DB.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("service-discovery"),
		Item:      sMap,
	})
	c.lock.Unlock()
	if err != nil {
		return errors.Wrap(err, "Failed to create/set serviceToHost map")
	}
	return nil
}

// Get gets value with key
func (c *Client) Get(key string) (string, error) {

	params := &dynamodb.QueryInput{
		KeyConditionExpression: aws.String("Service = :service"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":service": {
				S: aws.String(key),
			},
		},
		TableName: aws.String("service-discovery"),
	}

	c.lock.Lock()
	resp, err := c.DB.Query(params)
	c.lock.Unlock()
	if err != nil {
		return "", errors.Wrap(err, "Failed to get hosts corresponding to the service")
	}
	var pairs []keyValuePair
	err = dynamodbattribute.UnmarshalListOfMaps(resp.Items, &pairs)
	if err != nil {
		return "", errors.Wrap(err, "Failed to unmarshal dynamo attribute")
	}
	var res []string
	for index := range pairs {
		res = append(res, pairs[index].Host)
	}
	json, _ := json.Marshal(res)
	return string(json), nil
}

// Delete deletes service & host combination
func (c *Client) Delete(key, value string) error {
	s := keyValuePair{
		Service: key,
		Host:    value,
	}
	sMap, err := dynamodbattribute.MarshalMap(s)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal serviceToHost map")
	}
	c.lock.Lock()
	_, err = c.DB.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String("service-discovery"),
		Key:       sMap,
	})
	c.lock.Unlock()
	if err != nil {
		return errors.Wrap(err, "Failed to delete service and host")
	}
	return nil
}
