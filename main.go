package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/awoo-detat/moon/bot"
)

func main() {
	log.Println("starting")
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN not found in environment variables")
	}
	auth := "Bot " + token

	moon := bot.New(auth)
	log.Println("created")
	moon.Start()

	defer moon.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}
