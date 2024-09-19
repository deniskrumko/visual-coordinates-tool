/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/deniskrumko/visual-coordinates-tool/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt, os.Kill,
	)

	defer cancel()

	if err := cmd.Execute(ctx); err != nil {
		log.Println(err)
	}
}
