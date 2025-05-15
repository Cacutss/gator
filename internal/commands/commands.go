package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	config "github.com/Cacutss/blog-aggregator/internal/config"
	database "github.com/Cacutss/blog-aggregator/internal/database"
	RSS "github.com/Cacutss/blog-aggregator/internal/rss"
	"github.com/google/uuid"
	pq "github.com/lib/pq"
	"os"
	"time"
)

func ProcessError(err error) error {
	if err == nil {
		return nil
	}
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505":
			// Unique constraint violation
			return fmt.Errorf("User is already following this feed")
		case "23503":
			// foreign key missing
			return fmt.Errorf("Feed or user does not exist")
		case "23502":
			// violation on not null values
			return fmt.Errorf("Missing information")
		}
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("User or feed does not exist")
	}
	return fmt.Errorf("unknown error :%w", err)
}

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
		return ProcessError(err)
	}
	configUser := config.User{
		Name: user.Name,
		ID:   user.ID,
	}
	s.Config.SetUser(configUser)
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
	user, err := s.Db.CreateUser(context.Background(), params)
	if err != nil {
		return ProcessError(err)
	}
	configUser := config.User{
		Name: user.Name,
		ID:   user.ID,
	}
	s.Config.User = configUser
	config.Write(*s.Config)
	fmt.Printf("User %s setted.\n", user.Name)
	return nil
}

func HandlerReset(s *config.State, cmd Command) error {
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		return ProcessError(err)
	}
	fmt.Println("Users table resetted.")
	os.Exit(0)
	return nil
}

func HandlerUsers(s *config.State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return ProcessError(err)
	}
	for _, v := range users {
		fmt.Print(v.Name)
		if v.Name == s.Config.User.Name {
			fmt.Print(" (current)")
		}
		fmt.Println("")
	}
	return nil
}

func HandlerAgg(s *config.State, cmd Command) error {
	feed, err := RSS.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return ProcessError(err)
	}
	fmt.Println(feed)
	return nil
}

func middleWareLoggedIn(handler func(s *config.State, cmd Command, user database.User) error) func(*config.State, Command) error {
	return func(s *config.State, cmd Command) error {
		user, err := s.Db.GetUserById(context.Background(), s.Config.User.ID)
		if err != nil {
			return ProcessError(err)
		}
		return handler(s, cmd, user)
	}
}

func HandlerAddfeed(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 3 {
		return fmt.Errorf("addfeed needs 2 arguments, got less than expected")
	}
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[1],
		Url:       cmd.Args[2],
		UserID:    user.ID,
	}
	feed, err := s.Db.CreateFeed(context.Background(), params)
	if err != nil {
		return ProcessError(err)
	}
	followparams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    s.Config.User.ID,
		FeedID:    feed.ID,
	}
	_, err = s.Db.CreateFeedFollow(context.Background(), followparams)
	if err != nil {
		return ProcessError(err)
	}
	return nil
}

func HandlerFeeds(s *config.State, cmd Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return ProcessError(err)
	}
	for _, v := range feeds {
		user, err := s.Db.GetUserById(context.Background(), v.UserID)
		if err != nil {
			return ProcessError(err)
		}
		fmt.Printf("%s %s %s\n", v.Name, v.Url, user.Name)
	}
	return nil
}

func HandlerFollow(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("follow needs an url to follow")
	}
	feed, err := s.Db.GetFeedByUrl(context.Background(), cmd.Args[1])
	if err != nil {
		return ProcessError(err)
	}
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feeds_follow, err := s.Db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return ProcessError(err)
	}
	fmt.Printf("%s %s\n", feeds_follow.FeedName, feeds_follow.UserName)
	return nil
}

func HandlerFollowing(s *config.State, cmd Command, user database.User) error {
	followed_feeds, err := s.Db.GetFollowedFeeds(context.Background(), user.ID)
	if err != nil {
		return ProcessError(err)
	}
	fmt.Printf("User %s is following:\n\n", user.Name)
	for _, f := range followed_feeds {
		fmt.Printf("%s\n", f.Name)
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
	Result.register("agg", HandlerAgg)
	Result.register("addfeed", middleWareLoggedIn(HandlerAddfeed))
	Result.register("feeds", HandlerFeeds)
	Result.register("follow", middleWareLoggedIn(HandlerFollow))
	Result.register("following", middleWareLoggedIn(HandlerFollowing))
	return Result
}
