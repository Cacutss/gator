package config

import (
	"encoding/json"
	"fmt"
	database "github.com/Cacutss/blog-aggregator/internal/database"
	"log"
	"os"
)

const (
	configpath = "/.gatorconfig.json"
)

type Config struct {
	User  string `json:"current_user_name"`
	Dburl string `json:"db_url"`
}

type State struct {
	Config *Config
	Db     *database.Queries
}

func (c *Config) SetUser(user string) error {
	c.User = user
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
		return fmt.Errorf("Error writing config")
	}
	if err := os.WriteFile(path+configpath, data, 0644); err != nil {
		return fmt.Errorf("Error writing config file.")
	}
	return nil
}

func Read() Config {
	conf := Config{}
	path, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Couldn't find user's home directory")
	}
	data, err := os.ReadFile(path + configpath)
	if err != nil {
		_, err := os.Create(path + configpath)
		if err != nil {
			log.Fatal("failure to create config file")
		}
	}
	if err := json.Unmarshal(data, &conf); err != nil {
		log.Fatal(err)
	}
	return conf
}
