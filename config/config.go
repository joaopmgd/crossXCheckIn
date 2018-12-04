package config

import (
	"flag"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"
)

// Config will setup the Endpoints, the sources that will be requested, Log and Memory Cache configuration
type Config struct {
	Endpoints map[string]Endpoint `yaml:"endpoints"`
	Logger    *logrus.Logger
	Users     []User             `yaml:"users"`
	Requests  map[string]Request `yaml:"requests"`
}

// Endpoints for the future Requests
type Endpoint struct {
	Url  string `yaml:"url"`
	Type string `yaml:"type"`
}

type User struct {
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
	Time     string `yaml:"time"`
	Workout  string `yaml:"workout"`
}

type Request struct {
	Header map[string]string `yaml:"header"`
	Body   map[string]string `yaml:"body"`
	Param  map[string]string `yaml:"param"`
}

func GetConfig() (*Config, error) {
	p, err := getConfigPath()
	if err != nil {
		return nil, err
	}
	return readConfig(p)
}

// getConfigPath gets the path for the configuration file from the env var of from the arguments
func getConfigPath() (string, error) {
	// cp := os.Getenv("CONFIG")
	cp := "resources/config.yaml"
	if args := flag.Args(); len(args) == 1 {
		cp = args[0]
	}
	if cp == "" {
		return "", errors.New("Config not provided as env var or argument, or wrong number arguments")
	}
	return cp, nil
}

// readConfig reads the config file and unmarshal to the config
func readConfig(path string) (*Config, error) {
	// Read the file
	var Fs = afero.NewOsFs()
	f, err := afero.ReadFile(Fs, path)
	if err != nil {
		return nil, errors.Wrapf(err, "Error when reading the config on %s", path)
	}

	// Unmarshal the yaml of the config
	c := &Config{}
	if err = yaml.Unmarshal(f, c); err != nil {
		return nil, errors.Wrapf(err, "Error when unmarshaling the config at %s", path)
	}
	c.Logger = logrus.New()
	return c, nil
}
