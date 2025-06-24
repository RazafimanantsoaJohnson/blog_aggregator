package commands

import (
	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/config"
)

type State struct {
	Config *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	List map[string]func(*State, Command) error
}
