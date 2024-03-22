package proxy

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Authn Contains a list of users
type Authn struct {
	Users []User `yaml:"users"`
}

// User Identifies an username
type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// ParseConfig read a configuration file in the path `location` and returns an Authn object
func ParseConfig(location *string) (*Authn, error) {
	data, err := os.ReadFile(*location)
	if err != nil {
		return nil, err
	}
	authn := Authn{}
	err = yaml.Unmarshal([]byte(data), &authn)
	if err != nil {
		return nil, err
	}
	return &authn, nil
}
