package redis

import (
	"fmt"
	"testing"

	"golang.org/x/text/language"

	"github.com/ThatsMrTalbot/i18n"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRedisStorage(t *testing.T) {
	t.Parallel()

	Convey("Given an empty memory store", t, func() {
		storage, err := Connect("127.0.0.1:6379", "", 0)
		So(err, ShouldBeNil)

		Convey("When the default language is set", func() {
			storage.SetDefaultLanguage(language.English)

			Convey("Then the default language should be correct", func() {
				lang, err := storage.DefaultLanguage()
				So(err, ShouldBeNil)
				So(lang.String(), ShouldEqual, language.English.String())
			})
		})

		Convey("When the supported languages is set", func() {
			storage.StoreSupportedLanguage(language.English)
			storage.StoreSupportedLanguage(language.Spanish)
			storage.StoreSupportedLanguage(language.French)
			storage.DeleteSupportedLanguage(language.French)

			Convey("Then the supported languages should be correct", func() {
				langs, err := storage.SupportedLanguages()
				So(err, ShouldBeNil)
				So(langs, ShouldHaveLength, 2)
				So(langs, ShouldContainLanguage, language.Spanish)
				So(langs, ShouldContainLanguage, language.English)
			})
		})

		Convey("When an item is added to the memory store", func() {

			expected := &i18n.Translation{
				Lang:  language.English,
				Key:   "SomeKey",
				Value: "SomeValue",
			}

			err := storage.Store(expected)
			So(err, ShouldBeNil)

			Convey("Then it should be accessable", func() {
				results, err := storage.GetAll()
				So(err, ShouldBeNil)
				So(results, ShouldHaveLength, 1)
				So(results[0], ShouldResemble, expected)
			})
		})

		Convey("When an item is added twice to the memory store", func() {

			expected := &i18n.Translation{
				Lang:  language.English,
				Key:   "SomeKey",
				Value: "SomeValue",
			}

			replacement := &i18n.Translation{
				Lang:  language.English,
				Key:   "SomeKey",
				Value: "SomeOtherValue",
			}

			err := storage.Store(expected)
			So(err, ShouldBeNil)

			err = storage.Store(replacement)
			So(err, ShouldBeNil)

			Convey("Then only one should exists", func() {
				results, err := storage.GetAll()
				So(err, ShouldBeNil)
				So(results, ShouldHaveLength, 1)
			})

			Convey("Then the value should be updated", func() {
				results, err := storage.GetAll()
				So(err, ShouldBeNil)
				So(results[0], ShouldResemble, replacement)
			})
		})
	})

	Convey("Given a populated storage", t, func() {
		storage, err := Connect("127.0.0.1:6379", "", 0)
		So(err, ShouldBeNil)

		err = storage.Store(&i18n.Translation{
			Lang:  language.English,
			Key:   "SomeKey",
			Value: "SomeValue",
		})

		So(err, ShouldBeNil)

		Convey("When an item is deleted from the storage", func() {
			storage.Delete(&i18n.Translation{
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
