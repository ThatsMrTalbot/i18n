package i18n

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ThatsMrTalbot/scaffold"

	"golang.org/x/net/context"
	"golang.org/x/text/language"
)

var (
	languages = newRequestLanguageMap()
)

// Matcher get parses languages from http requests
type Matcher struct {
	i18n *I18n
}

// NewMatcher creates a new matcher
func NewMatcher(i18n *I18n) *Matcher {
	return &Matcher{
		i18n: i18n,
	}
}

func (matcher *Matcher) stripPrefix(prefix string, r *http.Request) {
	if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
		r.URL.Path = p
	}
}

// Middleware returns scaffold.Middleware that stores the language in the context
// redirecting if necessary
func (matcher *Matcher) Middleware() scaffold.Middleware {
	return scaffold.Middleware(func(next scaffold.Handler) scaffold.Handler {
		return scaffold.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			if tag, redirect := matcher.handle(w, r); !redirect {
				ctx = NewLanguageContext(ctx, tag)
				prefix := fmt.Sprintf("/%s", matcher.tagString(tag))
				matcher.stripPrefix(prefix, r)
				next.CtxServeHTTP(ctx, w, r)
			}
		})
	})
}

// NewLanguageContext stores a language tag in the context
func NewLanguageContext(ctx context.Context, tag language.Tag) context.Context {
	return context.WithValue(ctx, "i18n_tag", tag)
}

// GetLanguageFromContext returns the language from the context
func GetLanguageFromContext(ctx context.Context) language.Tag {
	if tag, ok := ctx.Value("i18n_tag").(language.Tag); ok {
		return tag
	}
	return language.Und
}

// Wrapper returns a http.Handler that stores languages in a sharded map for
// retrieval using GetLanguageFromRequest
func (matcher *Matcher) Wrapper(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tag, redirect := matcher.handle(w, r); !redirect {
			languages.Add(r, tag)
			defer languages.Delete(r)

			if handler != nil {
				prefix := fmt.Sprintf("/%s", matcher.tagString(tag))
				matcher.stripPrefix(prefix, r)
				handler.ServeHTTP(w, r)
			}
			http.NotFound(w, r)
		}
	})
}

// GetLanguageFromRequest gets the language from the internal map
func GetLanguageFromRequest(r *http.Request) language.Tag {
	if tag, ok := languages.Get(r); ok {
		return tag
	}
	return language.Und
}

func (matcher *Matcher) match(lang string) (matched language.Tag, match bool, valid bool, exact bool) {
	tag, err := language.Parse(lang)
	if err != nil || lang == "" {
		return language.Und, false, false, false
	}

	matched = language.Und
	exact = true
	match = false
	valid = true

	for {
		for _, supported := range matcher.i18n.GetSupportedLanguages() {
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

func (matcher *Matcher) handle(w http.ResponseWriter, r *http.Request) (language.Tag, bool) {
	path := r.URL.Path
	segments := strings.Split(path, "/")

	if len(segments) <= 1 {
		segments = []string{"", ""}
	}

	tag, match, valid, exact := matcher.match(segments[1])

	if !match {
		tags, _, _ := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
		tag = matcher.i18n.GetDefaultLanguage()

		for _, t := range tags {
			a, b, _, c := matcher.match(t.String())
			if b {
				tag, match, exact = a, b, c
				break
			}
		}
	}

	if !valid {
		str := matcher.tagString(tag)
		segments = append([]string{"", str}, segments[1:]...)
	}

	if !exact {
		str := matcher.tagString(tag)
		segments[1] = str
	}

	if !match || !valid || !exact {
		r.URL.Path = strings.Join(segments, "/")
		http.Redirect(w, r, r.URL.String(), 302)
		return tag, true
	}

	return tag, false
}

func (matcher *Matcher) tagString(tag language.Tag) string {
	str := tag.String()
	if str == "und" {
		return ""
	}
	return str
}
