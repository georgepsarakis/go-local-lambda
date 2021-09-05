package main

import (
	"github.com/georgepsarakis/go-local-lambda/lambda_api/configuration"
	"github.com/georgepsarakis/go-local-lambda/lambda_api/server"
	"log"
)

func main() {
	conf, err := configuration.LoadConfigurationFromFile("local_lambda.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = server.Start(conf)
	if err != nil {
		log.Fatal(err)
	}
}
