package i18n

import (
	"sync"

	"golang.org/x/text/language"
)

// Storage interface
type Storage interface {
	GetAll() ([]*Translation, error)
	Store(*Translation) error
	Delete(*Translation) error

	DefaultLanguage() (language.Tag, error)
	SupportedLanguages() ([]language.Tag, error)

	SetDefaultLanguage(language.Tag) error
	StoreSupportedLanguage(language.Tag) error
	DeleteSupportedLanguage(language.Tag) error
}

type inMemoryStorage struct {
	lock         sync.RWMutex
	translations []*Translation

	defaultLang    language.Tag
	supportedLangs []language.Tag
}

// NewInMemoryStorage Creates a non persistent in memory translation store
func NewInMemoryStorage() Storage {
	return new(inMemoryStorage)
}

func (storage *inMemoryStorage) SupportedLanguages() ([]language.Tag, error) {
	storage.lock.RLock()
	defer storage.lock.RUnlock()

	return storage.supportedLangs, nil
}

func (storage *inMemoryStorage) DefaultLanguage() (language.Tag, error) {
	storage.lock.RLock()
	defer storage.lock.RUnlock()

	return storage.defaultLang, nil
}

func (storage *inMemoryStorage) StoreSupportedLanguage(tag language.Tag) error {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	for _, l := range storage.supportedLangs {
		if l.String() == tag.String() {
			return nil
		}
	}

	storage.supportedLangs = append(storage.supportedLangs, tag)

	return nil
}

func (storage *inMemoryStorage) DeleteSupportedLanguage(tag language.Tag) error {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	for i, l := range storage.supportedLangs {
		if l.String() == tag.String() {
			storage.supportedLangs = append(storage.supportedLangs[:i], storage.supportedLangs[i+1:]...)
			return nil
		}
	}

	return nil
}

func (storage *inMemoryStorage) SetDefaultLanguage(tag language.Tag) error {
	err := storage.StoreSupportedLanguage(tag)
	if err != nil {
		return err
	}

	storage.lock.Lock()
	defer storage.lock.Unlock()

	storage.defaultLang = tag

	return nil
}

func (storage *inMemoryStorage) GetAll() ([]*Translation, error) {
	storage.lock.RLock()
	defer storage.lock.RUnlock()

	return storage.translations, nil
}

func (storage *inMemoryStorage) Store(translation *Translation) error {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	for _, t := range storage.translations {
		sameLang := t.Lang.String() == translation.Lang.String()
		sameKey := t.Key == translation.Key
		if sameLang && sameKey {
			t.Value = translation.Value
			return nil
		}
	}

	storage.translations = append(storage.translations, translation)

	return nil
}

func (storage *inMemoryStorage) Delete(translation *Translation) error {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	for i, t := range storage.translations {
		sameLang := t.Lang.String() == translation.Lang.String()
		sameKey := t.Key == translation.Key
		if sameLang && sameKey {
			storage.translations, storage.translations[len(storage.translations)-1] = append(storage.translations[:i], storage.translations[i+1:]...), nil
			return nil
		}
	}

	return nil
}
