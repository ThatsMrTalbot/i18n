package redis

import (
	"encoding/json"

	"github.com/ThatsMrTalbot/i18n"
	"golang.org/x/text/language"
)

type translationObject struct {
	Lang  string `json:"lang"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func encode(t *i18n.Translation) string {
	data, _ := json.Marshal(&translationObject{
		Lang:  t.Lang.String(),
		Key:   t.Key,
		Value: t.Value,
	})
	return string(data)
}

func decode(t string) (*i18n.Translation, error) {
	var obj translationObject
	err := json.Unmarshal([]byte(t), &obj)
	return &i18n.Translation{
		Lang:  language.Make(obj.Lang),
		Key:   obj.Key,
		Value: obj.Value,
	}, err
}
