package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"go-bot/lib/e"
	"go-bot/storage"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/bytedance/sonic/decoder"
)

const (
	defoultPermission = 0774
)

var (
	ErrNoFiles = errors.New("no files")
)

type Storage struct {
	basePath string
}

func NewStorage(basePath string) *Storage {
	return &Storage{basePath: basePath}
}

func (s *Storage) Save(page *storage.Page) error {
	fpath := filepath.Join(s.basePath, page.UserName)
	if err := os.MkdirAll(fpath, defoultPermission); err != nil {
		return e.Wrap(err, "failed to create directory")
	}

	fName, err := fileName(page)
	if err != nil {
		return e.Wrap(err, "failed to get file name")
	}

	fpath = filepath.Join(fpath, fName)

	file, err := os.Create(fpath)
	if err != nil {
		return e.Wrap(err, "failed to create file")
	}

	defer file.Close()

	if gob.NewEncoder(file).Encode(page); err != nil {
		return e.Wrap(err, "failed to encode page")
	}

	return nil
}

func (s *Storage) PickRandom(userName string) (*storage.Page, error) {
	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, e.Wrap(err, "failed to read directory")
	}

	if len(files) == 0 {
		return nil, ErrNoFiles
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s *Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap(err, "failed to get file name")

	}

	filePath := filepath.Join(s.basePath, p.UserName, fileName)
	if err := os.Remove(filePath); err != nil {
		return e.Wrap(err, fmt.Sprintf("failed to remove file %s", filePath))
	}

	return nil
}

func (s *Storage) IsExist(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap(err, "failed to get file name")
	}

	filePath := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err := os.Stat(filePath) {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, e.Wrap(err, fmt.Sprintf("failed to get file %s", filePath))
	}

	return true, nil
}

func (s *Storage) decodePage(filepath string) (*storage.Page, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, e.Wrap(err, "failed to open file")
	}

	file.Close()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap(err, "failed to decode page")
	}
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
