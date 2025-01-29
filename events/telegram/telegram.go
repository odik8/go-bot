package telegram

import (
	"errors"
	"go-bot/clients/telegram"
	"go-bot/events"
	"go-bot/lib/e"
	"go-bot/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var ErrorUnknownEventType = errors.New("unknown event type")
var ErrorUnknownMetaType = errors.New("unknown meta type")

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fecth(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap(err, "can't get events")
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))
	for _, u := range updates {
		res = append(res, evet(u))
	}

	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap(ErrorUnknownEventType, "can't process message")
	}
}

func (p *Processor) processMessage(event events.Event) error{
	meta, err := meta(event)
	if err != nil {
		return e.Wrap(err, "can't process message")
	}

	if err:=p.doCmd(event.Text, meta.ChatID, meta.Username); err!=nil {
		return e.Wrap(err, "can't process message")
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok == false {
		return Meta{}, e.Wrap(ErrorUnknownMetaType, "can't get meta")
	}

	return res, nil
}

func evet(update telegram.Update) events.Event {
	updateType := fetchType(update)
	res := events.Event{
		Type: updateType,
		Text: fetchText(update),
	}

	if updateType == events.Message {
		res.Meta = Meta{
			ChatID:   update.Message.Chat.ID,
			Username: update.Message.From.Username,
		}
	}

	return res
}

func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}

func fetchType(update telegram.Update) events.Type {
	if update.Message == nil {
		return events.Unkonwn
	}

	return events.Message
}
