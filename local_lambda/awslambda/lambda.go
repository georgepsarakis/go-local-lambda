package awslambda

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

type ConfigurationProvider interface {
	GetEnvironmentVariables(functionName string) (map[string]string, error)
}

type AWS struct {
	lambdaClient lambdaiface.LambdaAPI
}

func NewAWSProvider(client lambdaiface.LambdaAPI) *AWS {
	return &AWS{
		lambdaClient: client,
	}
}

func (a *AWS) GetEnvironmentVariables(functionName string) (map[string]string, error) {
	out, err := a.lambdaClient.GetFunctionConfiguration(&lambda.GetFunctionConfigurationInput{
		FunctionName: aws.String(functionName),
	})
	if err != nil {
		return nil, err
	}
	env := make(map[string]string, len(out.Environment.Variables))
	for key, value := range out.Environment.Variables {
		env[key] = *value
	}
	return env, nil
}
