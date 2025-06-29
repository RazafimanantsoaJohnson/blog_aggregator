package commands

import (
	"context"
	"fmt"
	"testing"
)

func TestGetRSSFeed(t *testing.T) {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Println(*feed)
	t.Errorf("form")
}
