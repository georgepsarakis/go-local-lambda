# go-local-lambda

Execute Go Lambda Functions locally using your AWS configuration environment variables.

## How it works

### RPC

Since Go is a compiled language, the [Handler](https://github.com/aws/aws-lambda-go) is wrapped in an RPC service. 
The Lambda service executes the uploaded executable (`main`) and the RPC server starts listening on a port defined in a reserved Environment Variable (`_LAMBDA_SERVER_PORT`).
Invocation requests are forwarded via an RPC client, and the Handler response is returned to the invoker.

The local-lambda server can accept Lambda Invoke requests from the AWS SDK by assigning its listening address (`http://localhost:8923` by default) as the Endpoint URL.

### Environment Variables

The Environment Variables for a Lambda function can be provisioned by different methods and may not be statically assignable.
Therefore, there might be a need to fetch the environment variables from the actual Lambda configuration.
The [GetFunctionConfiguration](https://docs.aws.amazon.com/lambda/latest/dg/API_GetFunctionConfiguration.html) can provide the current Lambda configuration.

The environment variables are automatically fetched from AWS when an executable is started by the local-lambda server.
Alternatively, for an external execution of the RPC server (e.g. for debugging) the sub-package
for configuration provisioning can be used to preconfigure the environment, in order to simulate the Lambda environment.
An additional custom variable `LOCAL_LAMBDA_ENDPOINT_URL` is also defined,
exposing the local-lambda listening address to the locally started services.
This URL can then be used in AWS SDK clients as the AWS Endpoint URL for Lambda-to-Lambda invocations.

Environment variables from your session (such as `HOME`, `AWS_ACCESS_KEY_ID`) are also added to each sub-process environment.

## Example Usage

The only required variables are the name of the function, as defined in AWS, and the listening port which must be unique across functions.
Defining the path to your Lambda function code file containing the `main` function will also start the service.

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

Omitting the `mainPath` setting will not start a sub-process, but incoming requests will still be forwarded to `localhost:<PORT>`. 

You can invoke your Lambda locally using the AWS CLI or any other client:

```
$ aws --no-sign-request --endpoint-url http://localhost:8923 \
  lambda invoke --function-name test-echo-lambda \
                --payload '{"hello": "world"}'
```

### Starting the Lambda RPC server externally

In order to define the port on which the Lambda RPC server will 
listen within the Go executable, you must simply include 
the `_LAMBDA_SERVER_PORT` environment variable in its environment variables.

```bash
$ go build main.go
$ _LAMBDA_SERVER_PORT=9001 ./main
```

```yaml
functions:
  -
    name: my-lambda
    port: 9001
```

## Credits

- https://github.com/djhworld/go-lambda-invoke