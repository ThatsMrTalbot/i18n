package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/ThatsMrTalbot/i18n"
	"golang.org/x/text/language"
)

type Storage struct {
	url string
}

func NewStorage(url string) *Storage {
	return &Storage{
		url: url,
	}
}

func decode(body []byte) ([]*i18n.Translation, error) {
	var objs []*translationObject
	err := json.Unmarshal(body, &objs)
	if err != nil {
		return nil, err
	}

	t := make([]*i18n.Translation, 0, len(objs))

	for _, obj := range objs {
		t = append(t, &i18n.Translation{
			Lang:  language.Make(obj.Lang),
			Key:   obj.Key,
			Value: obj.Value,
		})
	}
	return t, nil
}

func (storage *Storage) GetAll() ([]*i18n.Translation, error) {
	resp, err := http.Get(storage.url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return decode(body)
}

func (storage *Storage) Store(t *i18n.Translation) error {
	return errors.New("Not implemented")
}

func (storage *Storage) Delete(t *i18n.Translation) error {
	return errors.New("Not implemented")
}
