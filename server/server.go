package main

import (
	"context"
	"fmt"
)

func main() {
	s := NewMysqlServer(&MysqlHandler{})
	if err := s.Start(context.Background()); err != nil {
		fmt.Printf("%v\n", err)
	}
}
