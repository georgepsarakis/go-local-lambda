package configuration

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"sort"
)

const LambdaRpcPortVariableName = "_LAMBDA_SERVER_PORT"
const EndpointUrlVariableName = "LOCAL_LAMBDA_ENDPOINT_URL"

type LambdaEntry struct {
	// The Lambda function name
	Name string `yaml:"name"`
	// Additional environment variables
	Env map[string]string `yaml:"env"`
	// The path to the executable file which starts the handler.
	// Leaving the path empty will not start a sub-process,
	// but still route requests to the given port.
	MainPath string `yaml:"mainPath"`
	// The port for incoming Lambda RPC requests
	Port uint16 `yaml:"port"`
}

func (e *LambdaEntry) Environment(awsRemoteEnv map[string]string, endpointAddress string) (env []string) {
	env = os.Environ()
	env = append(env, fmt.Sprintf("%s=%d", LambdaRpcPortVariableName, e.Port))
	env = append(env, fmt.Sprintf("%s=http://%s", EndpointUrlVariableName, endpointAddress))
	for name, value := range awsRemoteEnv {
		env = append(env, fmt.Sprintf("%s=%s", name, value))
	}
	for name, value := range e.Env {
		env = append(env, fmt.Sprintf("%s=%s", name, value))
	}

	sort.Strings(env)
	return env
}

type Configuration struct {
	EndpointAddress string            `yaml:"endpointAddress"`
	Functions       []LambdaEntry     `yaml:"functions"`
	portByName      map[string]uint16 `yaml:"-"`
}

func (c *Configuration) FindPort(functionName string) uint16 {
	if c.portByName != nil {
		return c.portByName[functionName]
	}
	c.portByName = make(map[string]uint16, len(c.Functions))
	for _, fn := range c.Functions {
		c.portByName[fn.Name] = fn.Port
	}
	return c.portByName[functionName]
}

func FromFile(path string) (*Configuration, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	conf := Configuration{}
	err = yaml.Unmarshal(f, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

