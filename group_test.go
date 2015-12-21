package i18n

import (
	"testing"

	"golang.org/x/text/language"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGroup(t *testing.T) {
	t.Parallel()

	Convey("Given an i18n group", t, func() {
		storage := NewInMemoryStorage()
		i18n := New(storage)
		group := i18n.Group("SomeKey")

		Convey("When a translation is added", func() {
			expected := &Translation{
				Lang:  language.English,
				Key:   "SomeKey",
				Value: "SomeValue",
			}

			err := group.Add(expected)
			So(err, ShouldBeNil)

			Convey("Then the translation should be accessable", func() {
				result := group.Get(language.English, "SomeKey")
				So(result, ShouldNotBeNil)
				So(result, ShouldResemble, expected)
			})

			Convey("Then the translation should be accessable with locle string", func() {
				result, err := group.GetWithLangString("en", "SomeKey")
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)
				So(result, ShouldResemble, expected)
			})

			Convey("Then the translation should be accessable with a child language", func() {
				result := group.Get(language.BritishEnglish, "SomeKey")
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)
				So(result, ShouldResemble, expected)
			})

			Convey("Then the translation should exist in the storage", func() {
				results, err := storage.GetAll()
				So(err, ShouldBeNil)
				So(results, ShouldHaveLength, 1)
				So(results[0], ShouldResemble, expected)
			})

			Convey("Then the translation should be accessable through the helper method", func() {
				result1 := group.T(language.English.String(), "SomeKey")
				result2 := group.T(language.English, "SomeKey")
				So(result1, ShouldEqual, expected.Value)
				So(result2, ShouldEqual, expected.Value)
			})
		})
	})

	Convey("Given a popluated group", t, func() {
		storage := NewInMemoryStorage()
		i18n := New(storage)
		group := i18n.Group("SomeKey")

		expected := &Translation{
			Lang:  language.English,
			Key:   "SomeKey",
			Value: "SomeValue",
		}

		err := group.Add(expected)
		So(err, ShouldBeNil)

		Convey("When an item is deleted", func() {
			err := group.Delete(expected)
			So(err, ShouldBeNil)

			Convey("Then item should not be accessable", func() {
				result := group.Get(language.English, "SomeKey")
				So(result, ShouldBeNil)
			})

			Convey("Then item should not exist in storage", func() {
				results, err := storage.GetAll()
				So(err, ShouldBeNil)
				So(results, ShouldHaveLength, 0)
			})
		})
	})
}
