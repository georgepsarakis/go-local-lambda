package lambda

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/georgepsarakis/go-local-lambda/lambda_api/configuration"
)

// ConfigurationVariables fetches the environment variables that are configured in AWS
func ConfigurationVariables(awsClient lambdaiface.LambdaAPI, name string) (map[string]string, error) {
	out, err := awsClient.GetFunctionConfiguration(&lambda.GetFunctionConfigurationInput{
		FunctionName: aws.String(name),
	})
	if err != nil {
		return nil, err
	}

	env := make(map[string]string)
	for key, value := range out.Environment.Variables {
		env[key] = *value
	}
	return env, nil
}

func FetchEnvironment(config *configuration.LocalLambdaConfiguration, awsClient lambdaiface.LambdaAPI) (envByFunction map[string]map[string]string, err error) {
	envByFunction = make(map[string]map[string]string, len(config.Functions))
	for _, f := range config.Functions {
		env, err := ConfigurationVariables(awsClient, f.Name)
		if err != nil {
			return nil, err
		}
		envByFunction[f.Name] = env
	}
	return envByFunction, nil
}