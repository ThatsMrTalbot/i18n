package server

import (
	"encoding/json"

	"golang.org/x/text/language"

	"github.com/ThatsMrTalbot/i18n"
)

type translationObject struct {
	Lang  string `json:"lang"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type payload struct {
	DefaultLanguage    string               `json:"default"`
	SupportedLanguages []string             `json:"supported"`
	Translations       []*translationObject `json:"translations"`
}

func encode(translations []*i18n.Translation, supported []language.Tag, defaultLang language.Tag) []byte {
	s := make([]string, 0, len(supported))

	for _, i := range supported {
		s = append(s, i.String())
	}

	objs := make([]*translationObject, 0, len(translations))

	for _, item := range translations {
		objs = append(objs, &translationObject{
			Lang:  item.Lang.String(),
			Key:   item.Key,
			Value: item.Value,
		})
	}

	p := &payload{
		DefaultLanguage:    defaultLang.String(),
		SupportedLanguages: s,
		Translations:       objs,
	}

	data, _ := json.Marshal(p)
	return data
}

func decode(data []byte) ([]*i18n.Translation, []language.Tag, language.Tag, error) {
	var p payload
	err := json.Unmarshal(data, &p)

	if err != nil {
		return nil, nil, language.Und, err
	}

	s := make([]language.Tag, 0, len(p.SupportedLanguages))

	for _, i := range p.SupportedLanguages {
		s = append(s, language.Make(i))
	}

	d := language.Make(p.DefaultLanguage)

	t := make([]*i18n.Translation, 0, len(p.Translations))
	for _, i := range p.Translations {
		t = append(t, &i18n.Translation{
			Lang:  language.Make(i.Lang),
			Key:   i.Key,
			Value: i.Value,
		})
	}

	return t, s, d, nil
}
