package i18n

import (
	"sync"
	"time"

	"golang.org/x/text/language"
)

// Translation object
type Translation struct {
	Lang  language.Tag
	Key   string
	Value string
}

// I18n defines basic translation get/set methods
type I18n struct {
	lock sync.RWMutex

	storage      []Storage
	translations *Cache

	quit chan struct{}
}

// New translation manager
func New(storage ...Storage) *I18n {
	return &I18n{
		storage:      storage,
		translations: new(Cache),
	}
}

// Sync translations with database
func (i18n *I18n) Sync() error {
	i18n.lock.Lock()
	defer i18n.lock.Unlock()

	i18n.translations.Clear()

	for _, storage := range i18n.storage {
		translations, err := storage.GetAll()

		if err != nil {
			return err
		}

		for _, translation := range translations {
			i18n.translations.Add(translation)
		}
	}

	return nil
}

// SetRefreshInterval sets the inerval to sync the translations, 0 means no sync
func (i18n *I18n) SetRefreshInterval(d time.Duration) {
	i18n.lock.Lock()
	defer i18n.lock.Unlock()

	if i18n.quit != nil {
		close(i18n.quit)
		i18n.quit = nil
	}

	if d <= 0 {
		return
	}

	refresh := time.NewTicker(d)
	i18n.quit = make(chan struct{})

	go func(refresh *time.Ticker, quit chan struct{}) {
		for {
			select {
			case <-refresh.C:
				i18n.Sync()
			case <-quit:
				refresh.Stop()
				return
			}
		}
	}(refresh, i18n.quit)
}

// T is a helper method to get translation by lang string or language tag
func (i18n *I18n) T(lang interface{}, key string) string {
	var translation *Translation

	switch lang.(type) {
	case string:
		translation, _ = i18n.GetWithLangString(lang.(string), key)
	case language.Tag:
		translation = i18n.Get(lang.(language.Tag), key)
	}

	if translation != nil {
		return translation.Value
	}

	return ""
}

// Close must be called before going out of scope to stop the refresh goroutine
func (i18n *I18n) Close() error {
	if i18n.quit != nil {
		close(i18n.quit)
		i18n.quit = nil
	}
	return nil
}

// GetWithLangString parses the lang string before lookip up the translation
func (i18n *I18n) GetWithLangString(lang string, key string) (*Translation, error) {
	tag, err := language.Parse(lang)

	if err != nil {
		return nil, err
	}

	return i18n.Get(tag, key), nil
}

// Get translation
func (i18n *I18n) Get(lang language.Tag, key string) *Translation {
	i18n.lock.RLock()
	defer i18n.lock.RUnlock()

	for {
		if t := i18n.translations.Get(lang, key); t != nil {
			return t
		}

		if lang.IsRoot() {
			break
		}
		lang = lang.Parent()
	}

	return nil
}

// Add translation
func (i18n *I18n) Add(translation *Translation) error {
	i18n.lock.Lock()
	defer i18n.lock.Unlock()

	for _, storage := range i18n.storage {
		if err := storage.Store(translation); err != nil {
			return err
		}
	}

	i18n.translations.Add(translation)

	return nil
}

// Delete translation
func (i18n *I18n) Delete(translation *Translation) error {
	i18n.lock.Lock()
	defer i18n.lock.Unlock()

	for _, storage := range i18n.storage {
		if err := storage.Delete(translation); err != nil {
			return err
		}
	}

	i18n.translations.Delete(translation)

	return nil
}
