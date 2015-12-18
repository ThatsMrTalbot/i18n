package i18n

import "sync"

// Storage interface
type Storage interface {
	GetAll() ([]*Translation, error)
	Store(*Translation) error
	Delete(*Translation) error
}

type inMemoryStorage struct {
	lock         sync.RWMutex
	translations []*Translation
}

// NewInMemoryStorage Creates a non persistent in memory translation store
func NewInMemoryStorage() Storage {
	return new(inMemoryStorage)
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
		sameLocale := t.Locale.String() == translation.Locale.String()
		sameKey := t.Key == translation.Key
		if sameLocale && sameKey {
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
		sameLocale := t.Locale.String() == translation.Locale.String()
		sameKey := t.Key == translation.Key
		if sameLocale && sameKey {
			storage.translations, storage.translations[len(storage.translations)-1] = append(storage.translations[:i], storage.translations[i+1:]...), nil
			return nil
		}
	}

	return nil
}
