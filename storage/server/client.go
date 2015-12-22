package server

import (
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"golang.org/x/text/language"

	"github.com/ThatsMrTalbot/i18n"
)

type RequestMiddleware func(r *http.Request)

type Storage struct {
	url string

	middleware []RequestMiddleware

	lock           sync.Mutex
	updated        time.Time
	translations   []*i18n.Translation
	supportedLangs []language.Tag
	defaultLang    language.Tag
}

func NewStorage(url string, middleware ...RequestMiddleware) *Storage {
	return &Storage{
		url:        url,
		middleware: middleware,
	}
}

func (storage *Storage) sync() error {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	now := time.Now()
	if now.Sub(storage.updated) > 1*time.Minute {
		client := &http.Client{}
		req, err := http.NewRequest("GET", storage.url, nil)
		if err != nil {
			return err
		}

		for _, middleware := range storage.middleware {
			middleware(req)
		}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		t, s, d, err := decode(body)
		if err != nil {
			return err
		}

		storage.translations = t
		storage.supportedLangs = s
		storage.defaultLang = d
		storage.updated = now
	}
	return nil
}

func (storage *Storage) GetAll() ([]*i18n.Translation, error) {
	err := storage.sync()
	if err != nil {
		return nil, err
	}

	return storage.translations, err
}

func (storage *Storage) Store(t *i18n.Translation) error {
	return errors.New("Not implemented")
}

func (storage *Storage) Delete(t *i18n.Translation) error {
	return errors.New("Not implemented")
}

func (storage *Storage) DefaultLanguage() (language.Tag, error) {
	err := storage.sync()
	if err != nil {
		return language.Und, err
	}

	return storage.defaultLang, err
}

func (storage *Storage) SupportedLanguages() ([]language.Tag, error) {
	err := storage.sync()
	if err != nil {
		return nil, err
	}

	return storage.supportedLangs, err
}

func (storage *Storage) SetDefaultLanguage(language.Tag) error {
	return errors.New("Not implemented")
}

func (storage *Storage) StoreSupportedLanguage(language.Tag) error {
	return errors.New("Not implemented")
}

func (storage *Storage) DeleteSupportedLanguage(language.Tag) error {
	return errors.New("Not implemented")
}
