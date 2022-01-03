package local_lambda

import (
	"bufio"
	"errors"
	"github.com/georgepsarakis/go-local-lambda/local_lambda/awslambda"
	"github.com/georgepsarakis/go-local-lambda/local_lambda/configuration"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"io"
	"os/exec"
)


type Subprocess struct {
	path string
	env []string
	cmd *exec.Cmd
}

func NewSubprocess(path string, env []string) *Subprocess {
	s := Subprocess{path: path, env: env}
	s.cmd = exec.Command("go", "run", s.path)
	s.cmd.Env = s.env
	return &s
}

func (s *Subprocess) Start() error {
	return s.cmd.Start()
}

func (s *Subprocess) Tail(logger *zap.Logger) (*errgroup.Group, error) {
	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	group := errgroup.Group{}
	tail := func(stream io.ReadCloser, typ string) error {
		logFields := []zap.Field{
			zap.String("stream", typ),
			zap.String("path", s.path),
		}
		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			logger.Info(scanner.Text(), logFields...)
		}
		return scanner.Err()
	}
	group.Go(func() error {
		return tail(stdout, "stdout")
	})
	group.Go(func() error {
		return tail(stderr, "stderr")
	})
	return &group, nil
}

func (s *Subprocess) Wait() error {
	return s.cmd.Wait()
}

func StartAll(logger *zap.Logger, config configuration.Configuration, provider awslambda.ConfigurationProvider) error {
	g := errgroup.Group{}
	for i := range config.Functions {
		f := config.Functions[i]
		if f.MainPath == "" {
			logger.Info("executable path not defined - skipping sub-process start", zap.String("function", f.Name))
			continue
		}
		g.Go(func() error {
			var err error
			remoteEnv, err := provider.GetEnvironmentVariables(f.Name)
			if err != nil {
				return errors.New("cannot resolve remote environment variables: "+err.Error())
			}
			env := f.Environment(remoteEnv, config.EndpointAddress)
			proc := NewSubprocess(f.MainPath, env)
			tail, err := proc.Tail(logger)
			if err != nil {
				return err
			}
			defer func() {
				err = tail.Wait()
			}()

			logger.Info("starting sub-process", zap.String("command", proc.cmd.String()))
			logger.Info("remote environment variables", zap.Any("aws", remoteEnv))
			if err := proc.Start(); err != nil {
				return err
			}
			err = proc.Wait()
			return err
		})
	}
	return g.Wait()
}
