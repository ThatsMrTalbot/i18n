package server

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"golang.org/x/text/language"

	"github.com/ThatsMrTalbot/i18n"
	. "github.com/smartystreets/goconvey/convey"
)

func TestServerStorage(t *testing.T) {
	t.Parallel()

	Convey("Given an empty memory store", t, func() {
		mem := i18n.NewInMemoryStorage()
		server := NewServer(mem)
		host := httptest.NewServer(server)

		storage := NewStorage(host.URL)

		Convey("When the backing default language is set", func() {
			mem.SetDefaultLanguage(language.English)

			Convey("Then the backing default language should be correct", func() {
				lang, err := storage.DefaultLanguage()
				So(err, ShouldBeNil)
				So(lang.String(), ShouldEqual, language.English.String())
			})
		})

		Convey("When the supported languages is set", func() {
			mem.StoreSupportedLanguage(language.English)
			mem.StoreSupportedLanguage(language.Spanish)
			mem.StoreSupportedLanguage(language.French)
			mem.DeleteSupportedLanguage(language.French)

			Convey("Then the supported languages should be correct", func() {
				langs, err := storage.SupportedLanguages()
				So(err, ShouldBeNil)
				So(langs, ShouldHaveLength, 2)
				So(langs, ShouldContainLanguage, language.Spanish)
				So(langs, ShouldContainLanguage, language.English)
			})
		})

		Convey("When an item is added to the backing memory store", func() {

			expected := &i18n.Translation{
				Lang:  language.English,
				Key:   "SomeKey",
				Value: "SomeValue",
			}

			err := mem.Store(expected)
			So(err, ShouldBeNil)

			Convey("Then it should be accessable", func() {
				results, err := storage.GetAll()
				So(err, ShouldBeNil)
				So(results, ShouldHaveLength, 1)
				So(results[0], ShouldResemble, expected)
			})
		})

		Reset(func() {
			host.Close()
		})
	})

	Convey("Given a populated storage", t, func() {
		mem := i18n.NewInMemoryStorage()
		server := NewServer(mem)
		host := httptest.NewServer(server)

		storage := NewStorage(host.URL)

		err := mem.Store(&i18n.Translation{
			Lang:  language.English,
			Key:   "SomeKey",
			Value: "SomeValue",
		})

		So(err, ShouldBeNil)

		Convey("When an item is deleted from the backing storage", func() {
			mem.Delete(&i18n.Translation{
				Lang:  language.English,
				Key:   "SomeKey",
				Value: "SomeValue",
			})

			Convey("The item should not exists in the storage", func() {
				results, err := storage.GetAll()
				So(err, ShouldBeNil)
				So(results, ShouldHaveLength, 0)
			})
		})

		Reset(func() {
			host.Close()
		})
	})
}

func ShouldContainLanguage(actual interface{}, expected ...interface{}) string {
	haystack, ok := actual.([]language.Tag)
	if !ok {
		return "This assertion requires the actual value to be of type []language.Tag"
	}

	if len(expected) != 1 {
		return "This assertion requires exactly 1 comparison values (you provided 0)."
	}

	needle, ok := expected[0].(language.Tag)

	if !ok {
		return "This assertion requires the comparison value to be of type language.Tag"
	}

	for _, item := range haystack {
		if item.String() == needle.String() {
			return ""
		}
	}

	return fmt.Sprintf("Expected collection to contain %s but it did not!", needle.String())
}
