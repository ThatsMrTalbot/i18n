package redis

import (
	"github.com/ThatsMrTalbot/i18n"
	"golang.org/x/text/language"
	"gopkg.in/redis.v3"
)

const (
	RedisKey                   = "i18n_translations"
	RedisDefaultLanguageKey    = "i18n_translations_default"
	RedisSupportedLanguagesKey = "i18n_translations_supported"
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

func (storage *Storage) SupportedLanguages() ([]language.Tag, error) {
	cmd := storage.client.LRange(RedisSupportedLanguagesKey, 0, -1)
	results, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	langs := make([]language.Tag, 0, len(results))

	for _, result := range results {
		langs = append(langs, language.Make(result))
	}

	return langs, nil
}

func (storage *Storage) DefaultLanguage() (language.Tag, error) {
	cmd := storage.client.Get(RedisDefaultLanguageKey)
	lang, err := cmd.Result()
	return language.Make(lang), err
}

func (storage *Storage) StoreSupportedLanguage(tag language.Tag) error {
	tx, err := storage.client.Watch(RedisSupportedLanguagesKey)
	if err != nil {
		return err
	}
	defer tx.Close()

	cmd := tx.LRange(RedisSupportedLanguagesKey, 0, -1)
	results, err := cmd.Result()
	if err != nil {
		return err
	}

	_, err = tx.Exec(func() error {
		for _, result := range results {
			if result == tag.String() {
				return nil
			}
		}
		tx.LPush(RedisSupportedLanguagesKey, tag.String())
		return nil
	})

	if err == redis.TxFailedErr {
		return storage.StoreSupportedLanguage(tag)
	}

	return err
}

func (storage *Storage) DeleteSupportedLanguage(tag language.Tag) error {
	tx, err := storage.client.Watch(RedisSupportedLanguagesKey)
	if err != nil {
		return err
	}
	defer tx.Close()

	cmd := tx.LRange(RedisSupportedLanguagesKey, 0, -1)
	results, err := cmd.Result()
	if err != nil {
		return err
	}

	_, err = tx.Exec(func() error {
		for i, result := range results {
			if result == tag.String() {
				tx.LSet(RedisSupportedLanguagesKey, int64(i), "~REMOVE~")
			}
		}
		tx.LRem(RedisSupportedLanguagesKey, 0, "~REMOVE~")
		return nil
	})

	if err == redis.TxFailedErr {
		return storage.StoreSupportedLanguage(tag)
	}

	return err
}

func (storage *Storage) SetDefaultLanguage(tag language.Tag) error {
	return storage.client.Set(RedisDefaultLanguageKey, tag.String(), 0).Err()
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
