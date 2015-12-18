package i18n

import (
	"net/http"
	"strings"
	"sync"

	"golang.org/x/text/language"
)

var (
	languages struct {
		sync.RWMutex
		values map[*http.Request]language.Tag
	}
)

func init() {
	languages.values = make(map[*http.Request]language.Tag)
}

func GetLanguage(r *http.Request) language.Tag {
	languages.RLock()
	defer languages.RUnlock()

	if tag, ok := languages.values[r]; ok {
		return tag
	}
	return language.Und
}

type Router struct {
	DefaultLanguage    language.Tag
	SupportedLanguages []language.Tag
	Handler            http.Handler
}

func (router *Router) match(locale string) (matched language.Tag, match bool, valid bool, exact bool) {
	tag, err := language.Parse(locale)
	if err != nil || locale == "" {
		return language.Und, false, false, false
	}

	matched = language.Und
	exact = true
	match = false
	valid = true

	for {
		for _, supported := range router.SupportedLanguages {
			if supported.String() == tag.String() {
				matched = supported
				match = true
				return
			}
		}

		exact = false

		if tag.IsRoot() {
			break
		}

		tag = tag.Parent()
	}

	return
}

func (router *Router) tagString(tag language.Tag) string {
	str := tag.String()
	if str == "und" {
		return ""
	}
	return str
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	segments := strings.Split(path, "/")

	if len(segments) <= 1 {
		segments = []string{"", ""}
	}

	tag, match, valid, exact := router.match(segments[1])

	if !match {
		//TODO: Look at what the browser accepts
		tag = router.DefaultLanguage
	}

	if !valid {
		segments = append([]string{"", tag.String()}, segments[1:]...)
	}

	if !exact {
		segments[1] = tag.String()
	}

	if !match || !valid || !exact {
		r.URL.Path = strings.Join(segments, "/")
		http.Redirect(w, r, r.URL.String(), 302)
		return
	}

	languages.Lock()
	languages.values[r] = tag
	languages.Unlock()

	defer func() {
		languages.Lock()
		delete(languages.values, r)
		languages.Unlock()
	}()

	if router.Handler == nil {
		router.Handler = http.NotFoundHandler()
	}

	http.StripPrefix(strings.Join(segments[:2], "/"), router.Handler).ServeHTTP(w, r)
}
