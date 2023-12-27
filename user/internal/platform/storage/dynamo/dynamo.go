package dynamo

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	user "github.com/mata649/cqrs_on_aws/user/internal"
)

type UserDynamoRepository struct {
	db        *dynamodb.DynamoDB
	tableName string
}

func NewUserDynamoRepository(tableName string) *UserDynamoRepository {
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		},
	))
	return &UserDynamoRepository{
		db:        dynamodb.New(sess),
		tableName: tableName,
	}
}

func (repository UserDynamoRepository) Close() {

}
func (repository UserDynamoRepository) Create(ctx context.Context, user user.User) error {
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(repository.tableName),
	}
	_, err = repository.db.PutItem(input)
	if err != nil {
		return err
	}
	return err

}
func (repository UserDynamoRepository) GetByID(ctx context.Context, id string) (user.User, error) {

	keyEx := expression.Key("id").Equal(expression.Value(id))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return user.User{}, err
	}

	result, err := repository.db.QueryWithContext(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(repository.tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		return user.User{}, err
	}

	users := []user.User{}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return user.User{}, err
	}
	if len(users) == 0 {
		return user.User{}, nil
	}
	return users[0], nil
}
func (repository UserDynamoRepository) GetByUsername(ctx context.Context, username string) (user.User, error) {
	keyEx := expression.Key("username").Equal(expression.Value(username))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return user.User{}, err
	}
	result, err := repository.db.QueryWithContext(ctx, &dynamodb.QueryInput{
		TableName:                 &repository.tableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		IndexName:                 aws.String("username_index"),
	})
	if err != nil {
		return user.User{}, err
	}

	users := []user.User{}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return user.User{}, err
	}
	if len(users) == 0 {
		return user.User{}, nil
	}
	return users[0], nil
}
func (repository UserDynamoRepository) Get(ctx context.Context) ([]user.User, error) {

	result, err := repository.db.Query(&dynamodb.QueryInput{})
	if err != nil {
		return nil, err
	}
	users := []user.User{}
	for _, item := range result.Items {
		user := user.User{}
		err = dynamodbattribute.UnmarshalMap(item, user)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}
func (repository UserDynamoRepository) Delete(ctx context.Context, id string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: &id,
			},
		},
		TableName: aws.String(repository.tableName),
	}
	_, err := repository.db.DeleteItem(input)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func (repository UserDynamoRepository) Update(ctx context.Context, user user.User) error {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":usrn": {
				S: aws.String(user.Username().String()),
			},
			":passw": {
				S: aws.String(user.Password().String()),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(user.ID().String()),
			},
		},
		TableName:        &repository.tableName,
		UpdateExpression: aws.String("set username = :usrn, password = :passw"),
	}
	_, err := repository.db.UpdateItem(input)
	if err != nil {
		return err
	}
	return nil
}
