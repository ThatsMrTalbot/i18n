package i18n

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/text/language"

	. "github.com/smartystreets/goconvey/convey"
)

type TestHandler struct {
	Tag     language.Tag
	Request *http.Request
}

func (h *TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Tag = GetLanguage(r)
	h.Request = r
}

func TestRouter(t *testing.T) {
	t.Parallel()

	Convey("Given a server handling language selection", t, func() {
		handler := &TestHandler{}

		r := &Router{
			DefaultLanguage: language.English,
			SupportedLanguages: []language.Tag{
				language.English,
				language.Spanish,
				language.BritishEnglish,
			},
			Handler: handler,
		}

		server := httptest.NewServer(r)

		Convey("When I go to url with valid lang", func() {
			url := fmt.Sprintf("%s/%s/%s", server.URL, language.Spanish.String(), "other")
			_, err := http.Get(url)
			So(err, ShouldBeNil)

			Convey("Then the correct lang and path should be forwarded to the handler", func() {
				So(handler.Tag, ShouldResemble, language.Spanish)
				So(handler.Request.URL.Path, ShouldEqual, "/other")
			})
		})

		Convey("When I go to url with valid child lang", func() {
			url := fmt.Sprintf("%s/%s/%s", server.URL, language.AmericanEnglish.String(), "other")
			_, err := http.Get(url)
			So(err, ShouldBeNil)

			Convey("Then the parent lang and correct path should be forwarded to the handler", func() {
				So(handler.Tag, ShouldResemble, language.English)
				So(handler.Request.URL.Path, ShouldEqual, "/other")
			})
		})

		Convey("When I go to url with invalid lang", func() {
			url := fmt.Sprintf("%s/%s/%s", server.URL, language.French.String(), "other")
			_, err := http.Get(url)
			So(err, ShouldBeNil)

			Convey("Then the default lang and correct path should be forwarded to the handler", func() {
				So(handler.Tag, ShouldResemble, r.DefaultLanguage)
				So(handler.Request.URL.Path, ShouldEqual, "/other")
			})
		})

		Convey("When I go to url with no lang", func() {
			url := fmt.Sprintf("%s/%s", server.URL, "other")
			_, err := http.Get(url)
			So(err, ShouldBeNil)

			Convey("Then the default lang and correct path should be forwarded to the handler", func() {
				So(handler.Tag, ShouldResemble, r.DefaultLanguage)
				So(handler.Request.URL.Path, ShouldEqual, "/other")
			})
		})

		Convey("When I finish any request", func() {
			url := fmt.Sprintf("%s/%s", server.URL, "other")
			_, err := http.Get(url)
			So(err, ShouldBeNil)

			Convey("Then request should be removed from the internal language map", func() {
				So(languages.values, ShouldHaveLength, 0)
			})
		})

		Reset(func() {
			server.Close()
		})
	})
}
