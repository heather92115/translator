package db

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// DbConnect holds our db connection info
type DbConnect struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	DbName   string `json:"dbname"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

// dbConnectFromJson Unmarshalls a JSON string into the DbConnect instance.
func dbConnectFromJson(jsonStr string) (*DbConnect, error) {

	// Convert string to byte slice
	jsonData := []byte(jsonStr)

	// Declare a variable of type DbConnect
	var dbConn DbConnect

	// Decode the JSON data into the struct
	err := json.Unmarshal(jsonData, &dbConn)
	if err != nil {
		fmt.Printf("error decoding JSON, %v", err)
		return nil, err
	}

	return &dbConn, nil
}

// lookupUrl Ask AWS to get us our db connection info
func lookupUrl(dbLink string, region string) (string, error) {

	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return "", err
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(sdkConfig)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(dbLink),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		return "", err
	}

	return *result.SecretString, nil

}

// getEnv retrieves environment variables or returns a default value
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func createUrl(dbConnect *DbConnect) string {
	// Construct the URL, ensuring special characters in the password are encoded
	password := url.QueryEscape(dbConnect.Password)
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConnect.Host,
		dbConnect.Port,
		dbConnect.UserName,
		password,
		dbConnect.DbName)

}

// GetDatabaseURL Get the database URL used to connect
func GetDatabaseURL() string {

	dbLink := getEnv("DB_LINK", "")
	if len(dbLink) == 0 {
		panic("No DB_LINK environment variable found with no remediation")
	}
	region := getEnv("REGION", "us-east-1")

	dbInfo, err := lookupUrl(dbLink, region)
	if err != nil {
		fmt.Printf("Failed to obtain database info, err %v", err)
		panic("Failed to obtain database info")
	}

	dbConnect, err := dbConnectFromJson(dbInfo)
	if err != nil {
		fmt.Printf("Failed to unmarshall db Info %s, err %v", dbInfo, err)
		panic("Failed to obtain database info")
	}

	return createUrl(dbConnect)
}
