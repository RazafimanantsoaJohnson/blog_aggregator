package main

import (
	"fmt"

	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Errorf(err.Error())
	}
	conf.SetUser("myUser")
	fmt.Println(config.Read())
}
