package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/georgepsarakis/go-local-lambda/lambda_api/configuration"
	lambda2 "github.com/georgepsarakis/go-local-lambda/lambda_api/lambda"
	"github.com/georgepsarakis/go-local-lambda/lambda_api/server"
	"github.com/georgepsarakis/go-local-lambda/lambda_api/subprocess"
	"log"
	"os"
)


func main() {
	os.Setenv("AWS_SDK_LOAD_CONFIG", "true")

	awsSession, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	awsLambdaClient := lambda.New(awsSession)

	configurationPath := "local_lambda.yml"
	if len(os.Args) > 1 {
		configurationPath = os.Args[1]
	}
	config, err := configuration.LoadConfigurationFromFile(configurationPath)
	if err != nil {
		log.Fatal(err)
	}
	if config.EndpointAddress == "" {
		config.EndpointAddress = server.DefaultListeningAddress
	}
	envByFunction, err := lambda2.FetchEnvironment(config, awsLambdaClient)
	if err != nil {
		log.Fatal(err)
	}
	err = subprocess.StartAll(config, envByFunction)
	if err != nil {
		log.Fatal(err)
	}
}
