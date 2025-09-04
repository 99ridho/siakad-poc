package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

var CurrentConfig Config

type DatabaseConfigParams struct {
	Hostname string `json:"hostname"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     string `json:"port"`
	Schema   string `json:"schema"`
}

func (c DatabaseConfigParams) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		c.Username,
		c.Password,
		c.Hostname,
		c.Port,
		c.Database,
	)
}

type JWTConfigParams struct {
	Secret string `json:"secret"`
}

type Config struct {
	Database DatabaseConfigParams `json:"database"`
	JWT      JWTConfigParams      `json:"jwt"`
}

func LoadConfig() error {
	file, err := os.ReadFile("./config.json")
	if err != nil {
		return errors.Wrap(err, "error loading config")
	}

	err = json.Unmarshal(file, &CurrentConfig)
	if err != nil {
		return errors.Wrap(err, "error loading config")
	}

	return nil
}
