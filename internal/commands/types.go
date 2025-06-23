package commands

import (
	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/config"
)

type State struct {
	config *config.Config
}

type Command struct {
	name string
	args []string
}

type Commands struct {
	list map[string]func(*State, Command) error
}
