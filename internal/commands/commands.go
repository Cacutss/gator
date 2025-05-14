package commands

import (
	"context"
	"fmt"
	config "github.com/Cacutss/blog-aggregator/internal/config"
	database "github.com/Cacutss/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"log"
	"os"
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
	_, err := s.Db.CreateUser(context.Background(), params)
	if err != nil {
		log.Fatal("User already exists")
	}
	HandlerLogin(s, cmd)
	return nil
}

func HandlerReset(s *config.State, cmd Command) error {
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Users table resetted.")
	os.Exit(0)
	return nil
}

func HandlerUsers(s *config.State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	for _, v := range users {
		fmt.Print(v.Name)
		if v.Name == s.Config.User {
			fmt.Print(" (current)")
		}
		fmt.Println("")
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
	Result.register("reset", HandlerReset)
	Result.register("users", HandlerUsers)
	return Result
}
