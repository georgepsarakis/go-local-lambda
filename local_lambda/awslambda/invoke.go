package awslambda

import (
	"fmt"
	"net/rpc"
	"time"

	"github.com/aws/aws-lambda-go/lambda/messages"
)

const RequestID = "79186900-3cb3-4b93-bc65-611d05663264"

type Request struct {
	Payload []byte
	Deadline time.Time
}

type remoteProcedureCaller interface {
	Call(serviceMethod string, args interface{}, reply interface{}) error
}

type Client struct {
	rpc          remoteProcedureCaller
	functionName string
}

func NewClient(functionName string, port uint16) (*Client, error) {
	rpcClient, err := rpc.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	return &Client{rpc: rpcClient, functionName: functionName}, nil
}

func (c *Client) Invoke(r Request) (*messages.InvokeResponse, error) {
	request, err := rpcRequest(c.functionName, r.Payload, r.Deadline)
	if err != nil {
		return nil, err
	}

	var response messages.InvokeResponse
	if err = c.rpc.Call("Function.Invoke", request, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func rpcRequest(functionName string, payload []byte, deadline time.Time) (*messages.InvokeRequest, error) {
	if deadline.IsZero() {
		deadline = time.Now().Add(900 * time.Second)
	}

	return &messages.InvokeRequest{
		Payload:      payload,
		RequestId:    RequestID,
		Deadline: messages.InvokeRequest_Timestamp{
			Seconds: deadline.Unix(),
		},
		InvokedFunctionArn: fmt.Sprintf("arn:aws:lambda:us-west-2:123456789012:function:%s", functionName),
	}, nil
}
