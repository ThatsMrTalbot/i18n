# i18n

[![GoDoc](https://godoc.org/github.com/ThatsMrTalbot/i18n?status.svg)](https://godoc.org/github.com/ThatsMrTalbot/i18n) [![Build Status](https://travis-ci.org/ThatsMrTalbot/i18n.svg)](https://travis-ci.org/ThatsMrTalbot/i18n) [![Coverage Status](https://coveralls.io/repos/ThatsMrTalbot/i18n/badge.svg?branch=master&service=github)](https://coveralls.io/github/ThatsMrTalbot/i18n?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/ThatsMrTalbot/i18n)](https://goreportcard.com/report/github.com/ThatsMrTalbot/i18n)

This package is a translation package for golang. It provides basic translation management and http routing.

```go
storage := i18n.NewInMemoryStorage() // this is non persistent storage, for testing only
t := i18n.New(storage)

value := t.GetWithLangString("en-GB", "SomeKey")
if value != nil {
    // Translation exists
}else{
    // Translation does not exist
}

// OR

valueString := t.T("en-GB", "SomeKey")

// OR

valueString := t.T(language.BritishEnglish, "SomeKey")

```

It allows background synchronization with the storage for updating translations.

```go
t := i18n.New(storage)
t.SetRefreshInterval(1 * time.Hour)
defer t.Close() // This must be called to stop the refresh goroutine
```

To use the http router you wrap your default router in the Router object. All URLs will be prefixed with the language code. A specific language will also be matched by a generic parent, so /en-GB/some/path will match en and be redirected to /en/some/path.

If no language is specified in the URL, or it is not supported the Accept-Language header will be used to determine language. If the Accept-Language header is not set then the default language will be used.

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
    language := GetLanguageFromRequest(r)

    value := handler.translations.T(language, "SomeKey")
    w.Write([]byte(value))
}

func main() {
    storage := i18n.NewInMemoryStorage()
    translations := i18n.New(storage)

    translations.Add(&i18n.Translation{
        Lang: language.English
        Key: "SomeKey",
        Value: "SomeValue"
    })

    defaultHandler := &SomeHandler{
        translations: translations,
    }

    matcher := NewMatcher(translations)

    http.ListenAndServe(":8080", matcher.Wrapper(defaultHandler))
}

```

You can also use it on conjunction with [scaffold](https://github.com/ThatsMrTalbot/scaffold). In this case the language is stored in the context.

```go
package main

import (
	"net/http"

    "github.com/ThatsMrTalbot/i18n"    
	"github.com/ThatsMrTalbot/scaffold"    
	"golang.org/x/text/language"
    "golang.org/x/net/context"
)

type SomeHandler struct {
    translations *i18n.I18n
}

func (handler *SomeHandler) CtxServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    language := GetLanguageFromContext(r)

    value := handler.translations.T(language, "SomeKey")
    w.Write([]byte(value))
}

func main() {
    storage := i18n.NewInMemoryStorage()
    translations := i18n.New(storage)

    translations.Add(&i18n.Translation{
        Lang: language.English
        Key: "SomeKey",
        Value: "SomeValue"
    })

    matcher := NewMatcher(translations)

    defaultHandler := &SomeHandler{
        translations: translations,
    }

    dispatcher := scaffold.DefaultDispatcher()
    router := scaffold.New(dispatcher)

    router.Use(matcher.Middleware())

    // Since all URLs will be prefixed by a language, you must take
    // that into consideration when routing
    router.Route(":lang").Handle("", defaultHandler)

    http.ListenAndServe(":8080", dispatcher)
}

```
