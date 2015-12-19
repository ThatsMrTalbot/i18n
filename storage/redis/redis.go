package redis

import (
	"encoding/json"

	"github.com/ThatsMrTalbot/i18n"
	"golang.org/x/text/language"
	"gopkg.in/redis.v3"
)

const (
	RedisKey = "i18n_translations"
)

type Storage struct {
	client *redis.Client
}

func New(client *redis.Client) *Storage {
	return &Storage{
		client: client,
	}
}

func Connect(addr string, password string, db int64) (*Storage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping().Result()

	if err != nil {
		return nil, err
	}

	return New(client), nil
}

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

func (storage *Storage) GetAll() ([]*i18n.Translation, error) {
	cmd := storage.client.LRange(RedisKey, 0, -1)
	results, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	translations := make([]*i18n.Translation, 0, len(results))

	for _, result := range results {
		translation, err := decode(result)
		if err != nil {
			return nil, err
		}

		translations = append(translations, translation)
	}

	return translations, err
}

func (storage *Storage) Store(t *i18n.Translation) error {
	tx, err := storage.client.Watch(RedisKey)
	if err != nil {
		return err
	}
	defer tx.Close()

	cmd := tx.LRange(RedisKey, 0, -1)
	results, err := cmd.Result()
	if err != nil {
		return err
	}

	_, err = tx.Exec(func() error {
		for i, result := range results {
			tr, err := decode(result)
			if err != nil {
				return err
			}

			if tr.Key == t.Key && tr.Lang.String() == t.Lang.String() {
				tx.LSet(RedisKey, int64(i), "~REMOVE~")
			}
		}
		tx.LRem(RedisKey, 0, "~REMOVE~")
		tx.LPush(RedisKey, encode(t))
		return nil
	})

	if err == redis.TxFailedErr {
		return storage.Store(t)
	}

	return err
}

func (storage *Storage) Delete(t *i18n.Translation) error {
	tx, err := storage.client.Watch(RedisKey)
	if err != nil {
		return err
	}
	defer tx.Close()

	cmd := tx.LRange(RedisKey, 0, -1)
	results, err := cmd.Result()
	if err != nil {
		return err
	}

	_, err = tx.Exec(func() error {
		for i, result := range results {
			tr, err := decode(result)
			if err != nil {
				return err
			}

			if tr.Key == t.Key && tr.Lang.String() == t.Lang.String() {
				tx.LSet(RedisKey, int64(i), "~REMOVE~")
			}
		}
		tx.LRem(RedisKey, 0, "~REMOVE~")
		return nil
	})

	if err == redis.TxFailedErr {
		return storage.Delete(t)
	}

	return err
}
