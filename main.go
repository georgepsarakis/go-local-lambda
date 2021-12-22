package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/georgepsarakis/go-local-lambda/local_lambda"
	"github.com/georgepsarakis/go-local-lambda/local_lambda/configuration"
	"github.com/georgepsarakis/go-local-lambda/local_lambda/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const defaultConfigurationPath = "local_lambda.yml"

func main() {
	logCfg := zap.NewProductionConfig()
	logCfg.Encoding = "console"
	logCfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	logger, err := logCfg.Build()
	if err != nil {
		panic(err)
	}
	configurationPath := defaultConfigurationPath
	if len(os.Args) > 1 {
		configurationPath = os.Args[1]
	}
	logger.Info("configuration file", zap.String("path", configurationPath))

	awsSession, err := session.NewSession()
	if err != nil {
		logger.Fatal("cannot initialize AWS session", zap.Error(err))
	}
	awsLambdaClient := lambda.New(awsSession)

	config, err := configuration.FromFile(configurationPath)
	if err != nil {
		logger.Fatal("cannot load configuration from file", zap.Error(err))
	}
	if config.EndpointAddress == "" {
		config.EndpointAddress = server.DefaultListeningAddress
	}
	c := local_lambda.New(logger, awsLambdaClient)
	if err := c.Run(*config); err != nil {
		logger.Fatal("cannot start Lambda proxy processes or HTTP server", zap.Error(err))
	}
}
