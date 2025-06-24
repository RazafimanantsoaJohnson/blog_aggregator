package commands

import "fmt"

func (c *Commands) Run(s *State, cmd Command) error {
	helper, ok := c.List[cmd.Name]
	if !ok {
		return fmt.Errorf("the specified command is not supported")
	}
	return helper(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.List[name] = f
}
