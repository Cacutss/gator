package commands

import (
	"context"
	"fmt"
	config "github.com/Cacutss/blog-aggregator/internal/config"
	database "github.com/Cacutss/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"log"
	"time"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handler map[string]func(*config.State, Command) error
}

func HandlerLogin(s *config.State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Error given no parameters")
	}
	user, err := s.Db.GetUser(context.Background(), cmd.Args[1])
	if err != nil {
		log.Fatal("User does not exist")
	}
	s.Config.SetUser(user.Name)
	config.Write(*s.Config)
	fmt.Printf("User %s setted.\n", user.Name)
	return nil
}

func HandlerRegister(s *config.State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Error given no parameters")
	}
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[1],
	}
	i, err := s.Db.CreateUser(context.Background(), params)
	if err != nil {
		log.Fatal("User already exists", err, i)
	}
	return nil
}

func (c *Commands) register(name string, handler func(*config.State, Command) error) {
	c.Handler[name] = handler
}

func GetCommands() Commands {
	Result := Commands{
		Handler: make(map[string]func(*config.State, Command) error),
	}
	Result.register("login", HandlerLogin)
	Result.register("register", HandlerRegister)
	return Result
}
