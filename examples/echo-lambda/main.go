package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)

func echoHandler(input string) (string, error) {
	fmt.Println(fmt.Sprintf("received input: %q", input))
	return fmt.Sprintf(`{"echo": %q}`, input), nil
}

func main() {
	lambda.Start(echoHandler)
}

