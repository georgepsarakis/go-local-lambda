package local_lambda

import (
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/georgepsarakis/go-local-lambda/local_lambda/awslambda"
	"github.com/georgepsarakis/go-local-lambda/local_lambda/configuration"
	"github.com/georgepsarakis/go-local-lambda/local_lambda/server"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Coordinator struct {
	lambdaClient lambdaiface.LambdaAPI
	logger *zap.Logger
}

func New(logger *zap.Logger, lambdaClient lambdaiface.LambdaAPI) *Coordinator {
	return &Coordinator{
		logger: logger,
		lambdaClient: lambdaClient,
	}
}

func (c *Coordinator) Run(conf configuration.Configuration) error {
	g := errgroup.Group{}
	g.Go(func() error {
		// Start the HTTP server
		return server.Start(c.logger, conf)
	})
	g.Go(func() error {
		// Start Lambda sub-processes
		return StartAll(c.logger, conf, awslambda.NewAWSProvider(c.lambdaClient))
	})
	return g.Wait()
}
