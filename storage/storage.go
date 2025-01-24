package storage

import (
	"crypto/sha1"
	"fmt"
	"go-bot/lib/e"
	"io"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExist(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p *Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap(err, "failed to write string to hash")
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap(err, "failed to write string to hash")
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
