package dynamo

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/mata649/cqrs_on_aws/domain/user"
)

type UserDynamoRepository struct {
	db        *dynamodb.DynamoDB
	tableName string
}
type UserD struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	CreatedAt    time.Time `json:"createdAt"`
	TasksCreated uint      `json:"tasksCreated"`
}

func NewUserDynamoRepository(tableName string, sess *session.Session) *UserDynamoRepository {

	return &UserDynamoRepository{
		db:        dynamodb.New(sess),
		tableName: tableName,
	}
}

func (repository UserDynamoRepository) Close() {

}

func (repository UserDynamoRepository) Create(ctx context.Context, user user.User) error {
	userD := UserD{
		ID:           user.ID().String(),
		Email:        user.Email().String(),
		CreatedAt:    user.CreatedAt().Time(),
		TasksCreated: user.TasksCreated().Uint(),
	}
	av, err := dynamodbattribute.MarshalMap(userD)
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
	result, err := repository.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(repository.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return user.User{}, err
	}
	if result.Item == nil {
		return user.User{}, nil
	}
	userD := UserD{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &userD)
	if err != nil {
		return user.User{}, err
	}
	userFound, err := user.NewUser(userD.ID, userD.Email, userD.CreatedAt, userD.TasksCreated)
	if err != nil {
		return user.User{}, err
	}
	return userFound, nil

}

func (repository UserDynamoRepository) GetByEmail(ctx context.Context, email string) (user.User, error) {
	keyEx := expression.Key("email").Equal(expression.Value(email))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return user.User{}, err
	}
	result, err := repository.db.QueryWithContext(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(repository.tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		IndexName:                 aws.String("email_index"),
	})
	if err != nil {
		return user.User{}, err
	}

	users := []UserD{}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return user.User{}, err
	}
	if len(users) == 0 {
		return user.User{}, nil
	}
	userFound, err := user.NewUser(users[0].ID, users[0].Email, users[0].CreatedAt, users[0].TasksCreated)
	if err != nil {
		return user.User{}, err
	}
	return userFound, nil
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
	userD := UserD{
		ID:           user.ID().String(),
		Email:        user.Email().String(),
		CreatedAt:    user.CreatedAt().Time(),
		TasksCreated: user.TasksCreated().Uint(),
	}
	av, err := dynamodbattribute.MarshalMap(userD)
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
