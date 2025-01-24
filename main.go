package main

import (
	"flag"
	"go-bot/clients/telegram"
	"log"
)

func main() {

	tgClient := telegram.NewClient(mustToken())

	// proccessor := processor.NewProcessor(tgClient)

	// fetcher := fetcher.NewFetcher(tgClient)

	// consumer.Start(fetcher, proccessor)
}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "Telegram bot token")

	flag.Parse()

	if *token == "" {
		log.Fatal("Token is required")
	}

	return *token
}
