# i18n

This package is a translation package for golang. It provides basic translation management and http routing.

```go
storage := i18n.NewInMemoryStorage() // this is non persistent storage, for testing only
t := i18n.New(storage)

value := t.GetWithLocaleString("en-GB", "SomeKey")
if value != nil {
    // Translation exists
}else{
    // Translation does not exist
}
```

It allows background synchronization with the storage for updating translations.

```go
t := i18n.New(storage)
t.SetRefreshInterval(1 * time.Hour)
defer t.Close() // This must be called to stop the refresh goroutine
```

To use the http router you wrap your default router in the Router object. All URLS will not be prefixed with the language code.

```go
package main

import (
	"net/http"

	"github.com/ThatsMrTalbot/i18n"    
	"golang.org/x/text/language"
)

type SomeHandler struct {
    translations *i18n.I18n
}

func (handler *SomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    language := GetLanguage(r)

    value := handler.translations.T(language, "SomeKey")
    w.Write([]byte(value))
}

func main() {
    storage := i18n.NewInMemoryStorage()
    translations := i18n.New(storage)

    translations.Add(&i18n.Translation{
        Locale: language.English
        Key: "SomeKey",
        Value: "SomeValue"
    })

    defaultHandler := &SomeHandler{
        translations: translations,
    }

    languageRouter := &i18n.Router{
        DefaultLanguage: language.English,
        SupportedLanguages: []language.Tag{
            language.English,
            language.Spanish,
            language.BritishEnglish,
        },
        Handler: defaultHandler,
    }

    http.ListenAndServe(":8080", languageRouter)
}

```
