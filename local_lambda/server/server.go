package server

import (
	"github.com/georgepsarakis/go-local-lambda/local_lambda/configuration"
	"github.com/georgepsarakis/go-local-lambda/local_lambda/server/handlers"
	"go.uber.org/zap"
	"net/http"
)

const DefaultListeningAddress = "127.0.0.1:8923"

func Start(logger *zap.Logger, config configuration.Configuration) error {
	listeningAddress := config.EndpointAddress
	if listeningAddress == "" {
		listeningAddress = DefaultListeningAddress
	}
	handler := handlers.InvokeHandler{Configuration: config, Logger: logger}
	http.HandleFunc("/", handler.Run)
	logger.Info("starting Lambda proxy server", zap.String("address", listeningAddress))
	return http.ListenAndServe(listeningAddress, nil)
}