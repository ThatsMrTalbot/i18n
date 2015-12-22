package i18n

import (
	"net/http"
	"strings"

	"golang.org/x/text/language"
)

var (
	languages *requestLanguageMap
)

func init() {
	languages = newRequestLanguageMap()
}

func GetLanguage(r *http.Request) language.Tag {
	if tag, ok := languages.Get(r); ok {
		return tag
	}
	return language.Und
}

type Router struct {
	i18n    *I18n
	handler http.Handler
}

func NewRouter(handler http.Handler, i18n *I18n) *Router {
	return &Router{
		handler: handler,
		i18n:    i18n,
	}
}

func (router *Router) match(lang string) (matched language.Tag, match bool, valid bool, exact bool) {
	tag, err := language.Parse(lang)
	if err != nil || lang == "" {
		return language.Und, false, false, false
	}

	matched = language.Und
	exact = true
	match = false
	valid = true

	for {
		for _, supported := range router.i18n.GetSupportedLanguages() {
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
		tags, _, _ := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
		tag = router.i18n.GetDefaultLanguage()

		for _, t := range tags {
			a, b, _, c := router.match(t.String())
			if b {
				tag, match, exact = a, b, c
				break
			}
		}
	}

	if !valid {
		str := router.tagString(tag)
		segments = append([]string{"", str}, segments[1:]...)
	}

	if !exact {
		str := router.tagString(tag)
		segments[1] = str
	}

	if !match || !valid || !exact {
		r.URL.Path = strings.Join(segments, "/")
		http.Redirect(w, r, r.URL.String(), 302)
		return
	}

	languages.Add(r, tag)
	defer languages.Delete(r)

	if router.handler == nil {
		router.handler = http.NotFoundHandler()
	}

	http.StripPrefix(strings.Join(segments[:2], "/"), router.handler).ServeHTTP(w, r)
}
