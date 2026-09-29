package main

import (
	"fmt"
	"os"

	"github.com/mustaphalimar/gopanel/internal/app"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "gopanel-api:%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	a, err := app.New()
	if err != nil {
		return err
	}

	defer a.Close()
	return a.RunAPI()
}
