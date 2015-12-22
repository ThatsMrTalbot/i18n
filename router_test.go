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
		i18n := New()
		r := NewRouter(handler, i18n)
		i18n.AddSupportedLanguage(language.English, language.Spanish, language.BritishEnglish, language.French)
		i18n.RemoveSupportedLanguage(language.French)
		i18n.SetDefaultLanguage(language.English)

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
				So(handler.Tag, ShouldResemble, i18n.GetDefaultLanguage())
				So(handler.Request.URL.Path, ShouldEqual, "/other")
			})
		})

		Convey("When I go to url with no lang", func() {
			url := fmt.Sprintf("%s/%s", server.URL, "other")
			_, err := http.Get(url)
			So(err, ShouldBeNil)

			Convey("Then the default lang and correct path should be forwarded to the handler", func() {
				So(handler.Tag, ShouldResemble, i18n.GetDefaultLanguage())
				So(handler.Request.URL.Path, ShouldEqual, "/other")
			})
		})

		Convey("When I go to url with no lang, but an Accept-Language header", func() {
			client := &http.Client{}
			url := fmt.Sprintf("%s/%s", server.URL, "other")
			req, err := http.NewRequest("GET", url, nil)
			So(err, ShouldBeNil)

			req.Header.Add("Accept-Language", "es-ES;q=0.8, es;q=0.7")
			_, err = client.Do(req)
			So(err, ShouldBeNil)

			Convey("The language provided in the header should be used", func() {
				So(handler.Tag, ShouldResemble, language.Spanish)
				So(handler.Request.URL.Path, ShouldEqual, "/other")
			})

		})

		Convey("When I finish any request", func() {
			client := &http.Client{}
			url := fmt.Sprintf("%s/%s", server.URL, "other")
			req, err := http.NewRequest("GET", url, nil)
			So(err, ShouldBeNil)

			_, err = client.Do(req)
			So(err, ShouldBeNil)

			Convey("Then request should be removed from the internal language map", func() {
				tag, ok := languages.Get(req)
				So(tag.String(), ShouldEqual, language.Und.String())
				So(ok, ShouldBeFalse)
			})
		})

		Reset(func() {
			server.Close()
		})
	})
}
