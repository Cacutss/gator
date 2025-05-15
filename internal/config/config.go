package config

import (
	"encoding/json"
	"fmt"
	database "github.com/Cacutss/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"os"
)

const (
	configpath = "/.gatorconfig.json"
)

type User struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}

type Config struct {
	User  User   `json:"current_user"`
	Dburl string `json:"db_url"`
}

type State struct {
	Config *Config
	User   database.User
	Db     *database.Queries
}

func (c *Config) SetUser(user User) error {
	c.User.Name = user.Name
	c.User.ID = user.ID
	err := Write(*c)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func Write(cfg Config) error {
	path, _ := os.UserHomeDir()
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := os.WriteFile(path+configpath, data, 0644); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func LoadConfig() (Config, error) {
	conf := Config{}
	path, err := os.UserHomeDir()
	if err != nil {
		return conf, fmt.Errorf("%w", err)
	}
	data, err := os.ReadFile(path + configpath)
	if os.IsNotExist(err) {
		if err := Write(conf); err != nil {
			return conf, fmt.Errorf("%w", err)
		}
	}
	if err := json.Unmarshal(data, &conf); err != nil {
		return conf, fmt.Errorf("%w", err)
	}
	return conf, nil
}
