# go-local-lambda

Execute Go Lambda Functions locally using your AWS configuration variables.

Define the paths to your Lambda code and the name of the function in AWS:

```yaml
functions:
    -
        # The name of the Lambda function in AWS, used to retrieve environment variables
        name: test-echo-lambda
        # Local function port, must be unique
        port: 9000
        # The path to the Lambda executable
        mainPath: examples/echo-lambda/main.go
```

`go-local-lambda` will fetch the environment variables from AWS, and will also add a variable `LOCAL_LAMBDA_ENDPOINT_URL`,
which can be used in AWS SDK clients for Lambda-to-Lambda invocations.

Environment variables from your session (such as `HOME`, `AWS_ACCESS_KEY_ID`) are also added to each sub-process environment.

You can invoke your Lambda locally using the AWS CLI or any other client:

```
$ aws --endpoint-url http://localhost:8923 lambda invoke --function-name test-echo-lambda  --payload '{"hello": "world"}'
```
