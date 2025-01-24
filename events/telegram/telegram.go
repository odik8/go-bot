package telegram

import "go-bot/clients/telegram"

type Processor struct {
	tg *telegram.Client
	offset int
	// storage storage.Storage
}
