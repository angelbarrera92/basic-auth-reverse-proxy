package proxy

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	envPrefix            = "BARP"
	randomPasswordLength = 24
	randomPasswordPool   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789"
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

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewAuthn() *Authn {
	return &Authn{}
}

// Merges the data from the provided object and returns self
func (a *Authn) Merge(other *Authn) error {
	a.Users = append(a.Users, other.Users...)

	return nil
}

func (a *Authn) Validate() error {
	if len(a.Users) == 0 {
		return fmt.Errorf("no authentication principals configured")
	}

	return nil
}

func (a *Authn) AddUser(username string, password string) error {
	if len(username) == 0 {
		return fmt.Errorf("username must not be empty")
	}

	if len(password) == 0 {
		password = randomPassword()

		log.Printf("Generating random password for %s: %s", username, password)
	}

	a.Users = append(a.Users, User{
		Username: username,
		Password: password,
	})

	return nil
}

func (a *Authn) ParseEnvironment() error {
	prefix := strings.Join([]string{envPrefix, "USERNAME", ""}, "_")

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key, value := pair[0], pair[1]
		postfix := strings.TrimPrefix(key, prefix)

		if postfix == key {
			continue
		}

		lookup := strings.Join([]string{envPrefix, "PASSWORD", postfix}, "_")
		username := strings.TrimSpace(value)
		password := os.Getenv(lookup)

		if err := a.AddUser(username, password); err != nil {
			return fmt.Errorf("invalid environment variable %s: %v", key, err)
		}
	}

	return nil
}

// ParseConfig read a configuration file in the path `location` and returns self
func (a *Authn) ParseFile(location string) error {
	data, err := ioutil.ReadFile(location)
	if err != nil {
		return err
	}

	authn := NewAuthn()
	err = yaml.Unmarshal([]byte(data), authn)
	if err != nil {
		return err
	}

	err = a.Merge(authn)
	if err != nil {
		return err
	}

	return nil
}

func randomPassword() string {
	var b strings.Builder

	pool := []rune(randomPasswordPool)

	for i := 0; i < randomPasswordLength; i++ {
		selection := rand.Intn(len(pool))
		b.WriteRune(pool[selection])
	}
	return b.String()
}
