package dynamo

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/mata649/cqrs_on_aws/domain/task"
)

type TaskDynamoRepository struct {
	db        *dynamodb.DynamoDB
	tableName string
}
type TaskD struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UserID      string    `json:"userID"`
}

func NewTaskDynamoRepository(tableName string, sess *session.Session) *TaskDynamoRepository {

	return &TaskDynamoRepository{
		db:        dynamodb.New(sess),
		tableName: tableName,
	}
}

func (repository TaskDynamoRepository) Close() error {
	return nil
}

func (repository TaskDynamoRepository) Create(ctx context.Context, task task.Task) error {
	taskD := TaskD{
		ID:          task.ID().String(),
		Title:       task.Title().String(),
		Description: task.Description().String(),
		CreatedAt:   task.CreatedAt().Time(),
		UserID:      task.UserID().String(),
	}
	av, err := dynamodbattribute.MarshalMap(taskD)
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
	return nil
}

func (repository TaskDynamoRepository) GetByID(ctx context.Context, id string) (task.Task, error) {
	result, err := repository.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(repository.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return task.Task{}, err
	}
	if result.Item == nil {
		return task.Task{}, nil
	}
	taskD := TaskD{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &taskD)
	if err != nil {
		return task.Task{}, err
	}
	taskFound, err := task.NewTask(taskD.ID, taskD.Title, taskD.Description, taskD.CreatedAt, taskD.UserID)
	if err != nil {
		return task.Task{}, err
	}
	return taskFound, nil
}

func (repository TaskDynamoRepository) Delete(ctx context.Context, id string) error {
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
		return err
	}
	return nil
}

func (repository TaskDynamoRepository) Update(ctx context.Context, task task.Task) error {
	taskD := TaskD{
		ID:          task.ID().String(),
		Title:       task.Title().String(),
		Description: task.Description().String(),
		CreatedAt:   task.CreatedAt().Time(),
		UserID:      task.UserID().String(),
	}
	av, err := dynamodbattribute.MarshalMap(taskD)
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
