package internal

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/joho/godotenv"
)

var awsSess *session.Session
var svc *dynamodb.DynamoDB
var err error

func init() {
	godotenv.Load(os.Getenv("OOPS_ENV_FILE"))

	if os.Getenv("DB_DRIVER") == "dynamo" {
		awsSess, err = session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
		if err != nil {
			panic("Unable to instantiate AWS session")
		}
		log.Println("Instantiated AWS Session")
		svc = dynamodb.New(awsSess)
	}

}

func AddDynamoSecret(secret string, ttl int64) (string, error) {

	b := make([]byte, 16)
	rand.Read(b)
	u := hex.EncodeToString(b)

	item := Secret{
		Secret:     secret,
		Expiration: ttl,
		OopsID:     u,
	}

	ddbItem, err := dynamodbattribute.MarshalMap(item)

	if err != nil {
		log.Println("Error creating DynamoDB item", err.Error())
		return "", err
	}

	input := &dynamodb.PutItemInput{
		Item:      ddbItem,
		TableName: aws.String(TableName),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		log.Println("Error adding item to DynamoDB", err.Error())
		return "", err
	}

	log.Println("Created secret", u)

	return u, nil

}

func ReturnDynamoSecret(id string) (string, error) {

	input := &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"OopsID": {
				S: aws.String(id),
			},
		},
	}

	result, err := svc.GetItem(input)

	if err != nil {
		log.Println("Error retrieving secret from Dynamo", err.Error())
		return "", nil
	}

	item := Secret{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)

	if err != nil {
		log.Println("Error unmarshalling DynamoDB respone", err.Error())
		return "", nil
	}

	if item.Secret == "" {
		return "Secret not found", nil
	}

	go deleteDynamoItemAfterView(id)

	return item.Secret, nil

}

func deleteDynamoItemAfterView(id string) error {

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"OopsID": {
				S: aws.String(id),
			},
		},
	}

	_, err := svc.DeleteItem(input)

	if err != nil {
		log.Println("Error deleting entry from DynamoDB", err.Error())
		return err
	}

	log.Println("Deleted secret", id)

	return nil

}
