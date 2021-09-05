package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type LocalLambdaDefinition struct {
	Name string `yaml:"name"`
	Env map[string]string `yaml:"env"`
	MainPath string `yaml:"mainPath"`
	Port uint16 `yaml:"port"`
}

type LocalLambdaConfiguration struct {
	EndpointAddress string            `yaml:"endpointAddress"`
	Functions []LocalLambdaDefinition `yaml:"functions"`
	NamePortMapping map[string]uint16 `yaml:"-"`
}

func (c *LocalLambdaConfiguration) FindPort(functionName string) uint16 {
	if c.NamePortMapping != nil {
		return c.NamePortMapping[functionName]
	}
	c.NamePortMapping = make(map[string]uint16, len(c.Functions))
	for _, fn := range c.Functions {
		c.NamePortMapping[fn.Name] = fn.Port
	}
	return c.NamePortMapping[functionName]
}

func LoadConfigurationFromFile(path string) (*LocalLambdaConfiguration, error) {
	conf := LocalLambdaConfiguration{}
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(f, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

