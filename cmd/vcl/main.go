package main

import (
	"log"

	"github.com/KeisukeYamashita/go-vcl/pkg/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
