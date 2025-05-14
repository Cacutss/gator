package main

import (
	"database/sql"
	commands "github.com/Cacutss/blog-aggregator/internal/commands"
	config "github.com/Cacutss/blog-aggregator/internal/config"
	database "github.com/Cacutss/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("error no command given")
	}
	state := config.State{}
	conf := config.Read()
	state.Config = &conf
	db, err := sql.Open("postgres", state.Config.Dburl)
	if err != nil {
		log.Fatal("error connecting to database:", err)
	}
	dbQueries := database.New(db)
	state.Db = dbQueries
	Commands := commands.GetCommands()
	command := commands.Command{
		Name: os.Args[1],
		Args: os.Args[1:],
	}
	if h, ok := Commands.Handler[os.Args[1]]; ok {
		if err := h(&state, command); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("No given parameters")
	}
}
