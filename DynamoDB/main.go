package Dynamo

import (
	"context"
	"fmt"
	"main/Schemas"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

const REGION = "us-east-1"

func QueryDB(item string, tableName string, indexName string, keyCondition string, expressionAttribute string) Schemas.LoanDynamo {
	//Create AWS Session
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = REGION
		return nil
	})
	if err != nil {
		panic(err)
	}
	dynamodbClient := dynamodb.NewFromConfig(cfg)

	//Create Query
	out, err := dynamodbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String(indexName),
		KeyConditionExpression: aws.String(keyCondition),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			expressionAttribute: &types.AttributeValueMemberS{Value: item},
		},
	})
	if err != nil {
		panic(err)
	}

	//Parsing item dynamoDB
	var result []Schemas.LoanDynamo
	for i := range out.Items {
		var object Schemas.LoanDynamo
		err = attributevalue.UnmarshalMap(out.Items[i], &object)
		if err != nil {
			fmt.Println("Error passing object: ", err)
		}
		if err != nil {
			fmt.Println("Error passing result: ", err)
		}
		result = append(result, object)
	}
	return result[0]
}

func GetClient(item string, tableName string, keyCondition string, expressionAttribute string) Schemas.ClientDynamo {
	//Create AWS Session
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = REGION
		return nil
	})
	if err != nil {
		panic(err)
	}
	dynamodbClient := dynamodb.NewFromConfig(cfg)

	out, err := dynamodbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String(keyCondition),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			expressionAttribute: &types.AttributeValueMemberS{Value: item},
		},
	})
	if err != nil {
		panic(err)
	}
	//Parsing item dynamoDB

	var result []Schemas.ClientDynamo
	for i := range out.Items {
		var object Schemas.ClientDynamo
		err = attributevalue.UnmarshalMap(out.Items[i], &object)
		if err != nil {
			fmt.Println("Error to passing object: ", err)
		}

		if err != nil {
			fmt.Println("Error to passing result: ", err)
		}
		result = append(result, object)

	}
	return result[0]
}

func QueryRUC(item string, tableName string, key string) Schemas.RUC {
	//Create AWS Session
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = REGION
		return nil
	})
	if err != nil {
		panic(err)
	}
	svc := dynamodb.NewFromConfig(cfg)

	//Select object
	out, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			key: &types.AttributeValueMemberS{Value: item},
		},
	})
	if err != nil {
		fmt.Println("Ruc not found")
	}
	//Parsing item dynamoDB
	var object Schemas.RUC
	err = attributevalue.UnmarshalMap(out.Item, &object)
	if err != nil {
		fmt.Println("Error parsing object: ", err)
	}
	//resultJson, err := json.MarshalIndent(object, "", " ")
	//if err != nil {
	//	fmt.Println("Error parsing result: ", err)
	//}
	//fmt.Println(string(resultJson))
	return object
}
