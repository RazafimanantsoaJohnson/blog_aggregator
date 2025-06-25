package commands

import (
	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/config"
	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/database"
)

type State struct {
	Config    *config.Config
	DbQueries *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	List map[string]func(*State, Command) error
}
