package server

import (
	"github.com/georgepsarakis/go-local-lambda/lambda_api/configuration"
	"github.com/georgepsarakis/go-local-lambda/lambda_api/handlers"
	"net/http"
)

const DefaultListeningAddress = "127.0.0.1:8923"

func Start(config *configuration.LocalLambdaConfiguration) error {
	handler := handlers.InvokeHandler{Configuration: config}
	http.HandleFunc("/", handler.Run)
	return http.ListenAndServe(config.EndpointAddress, nil)
}