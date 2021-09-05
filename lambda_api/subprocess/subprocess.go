package subprocess

import (
	"bufio"
	"fmt"
	"github.com/georgepsarakis/go-local-lambda/lambda_api/configuration"
	"github.com/georgepsarakis/go-local-lambda/lambda_api/server"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"os/exec"
	"sort"
	"strconv"
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

func (s *Subprocess) Tail() (*errgroup.Group, error) {
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
		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			fmt.Println(fmt.Sprintf("[%s] [%s] %s", s.path, typ, scanner.Text()))
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

func StartAll(config *configuration.LocalLambdaConfiguration, awsEnvByFunction map[string]map[string]string) error {
	g := errgroup.Group{}

	g.Go(func() error {
		return server.Start(config)
	})

	for _, f := range config.Functions {
		g.Go(func() error {
			proc := NewSubprocess(f.MainPath, resolveEnvironmentVariables(&f, awsEnvByFunction[f.Name], config.EndpointAddress))
			tail, err := proc.Tail()
			if err != nil {
				return err
			}
			defer tail.Wait()

			err = proc.Start()
			fmt.Println("Starting sub-process:", proc.cmd.String(), "additional environment variables:", extraEnvironmentVariables(proc.cmd.Env))
			if err != nil {
				return err
			}

			return proc.Wait()
		})
	}
	return g.Wait()
}

func extraEnvironmentVariables(all []string) []string {
	var extra []string
	for _, v := range all {
		found := false
		for _, s := range os.Environ() {
			if s == v {
				found = true
				break
			}
		}
		if !found {
			extra = append(extra, v)
		}
	}
	sort.Strings(extra)
	return extra
}

func resolveEnvironmentVariables(definition *configuration.LocalLambdaDefinition, awsEnv map[string]string, endpointAddress string) (env []string) {
	env = os.Environ()
	env = append(env, fmt.Sprintf("_LAMBDA_SERVER_PORT=%d", definition.Port))
	env = append(env, fmt.Sprintf("LOCAL_LAMBDA_ENDPOINT_URL=http://%s", endpointAddress))
	for name, value := range awsEnv {
		env = append(env, fmt.Sprintf("%s=%s", name, strconv.Quote(value)))
	}
	for name, value := range definition.Env {
		env = append(env, fmt.Sprintf("%s=%s", name, strconv.Quote(value)))
	}

	sort.Strings(env)
	return env
}
