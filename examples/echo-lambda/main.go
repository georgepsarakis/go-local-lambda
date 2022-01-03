package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)

func echoHandler(input map[string]string) (map[string]string, error) {
	fmt.Println(fmt.Sprintf("received input: %q", input))
	return input, nil
}

func main() {
	lambda.Start(echoHandler)
}

