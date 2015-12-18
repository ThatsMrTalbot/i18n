package i18n

import (
	"testing"
	"time"

	"golang.org/x/text/language"

	. "github.com/smartystreets/goconvey/convey"
)

func TestI18n(t *testing.T) {
	t.Parallel()

	Convey("Given an empty memory store", t, func() {
		storage := NewInMemoryStorage()
		i18n := New(storage)

		Convey("When a translation is added", func() {
			expected := &Translation{
				Lang: language.English,
				Key:    "SomeKey",
				Value:  "SomeValue",
			}

			err := i18n.Add(expected)
			So(err, ShouldBeNil)

			Convey("Then the translation should be accessable", func() {
				result := i18n.Get(language.English, "SomeKey")
				So(result, ShouldNotBeNil)
				So(result, ShouldResemble, expected)
			})

			Convey("Then the translation should be accessable with locle string", func() {
				result, err := i18n.GetWithLangString("en", "SomeKey")
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)
				So(result, ShouldResemble, expected)
			})

			Convey("Then the translation should be accessable with a child language", func() {
				result := i18n.Get(language.BritishEnglish, "SomeKey")
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
				result1 := i18n.T(language.English.String(), "SomeKey")
				result2 := i18n.T(language.English, "SomeKey")
				So(result1, ShouldEqual, expected.Value)
				So(result2, ShouldEqual, expected.Value)
			})
		})

		Convey("When a translation is added to the storage", func() {

			expected := &Translation{
				Lang: language.English,
				Key:    "SomeKey",
				Value:  "SomeValue",
			}

			err := storage.Store(expected)
			So(err, ShouldBeNil)

			Convey("Then a call to sync should add it to the translation list", func() {
				err := i18n.Sync()
				So(err, ShouldBeNil)

				result := i18n.Get(language.English, "SomeKey")
				So(result, ShouldNotBeNil)
				So(result, ShouldResemble, expected)
			})

			Convey("Then a refresh interval should add it to the translation list", func() {
				i18n.SetRefreshInterval(50 * time.Millisecond)
				time.Sleep(100 * time.Millisecond)

				result := i18n.Get(language.English, "SomeKey")
				So(result, ShouldNotBeNil)
				So(result, ShouldResemble, expected)

				// Test killing the refresh goroutine
				i18n.SetRefreshInterval(0)

				storage.Delete(expected)
				time.Sleep(100 * time.Millisecond)

				result = i18n.Get(language.English, "SomeKey")
				So(result, ShouldNotBeNil)
				So(result, ShouldResemble, expected)
			})
		})

		Reset(func() {
			i18n.Close()
		})
	})

	Convey("Given a popluated memory store", t, func() {
		storage := NewInMemoryStorage()
		i18n := New(storage)

		expected := &Translation{
			Lang: language.English,
			Key:    "SomeKey",
			Value:  "SomeValue",
		}

		err := i18n.Add(expected)
		So(err, ShouldBeNil)

		Convey("When an item is deleted", func() {
			err := i18n.Delete(expected)
			So(err, ShouldBeNil)

			Convey("Then item should not be accessable", func() {
				result := i18n.Get(language.English, "SomeKey")
				So(result, ShouldBeNil)
			})

			Convey("Then item should not exist in storage", func() {
				results, err := storage.GetAll()
				So(err, ShouldBeNil)
				So(results, ShouldHaveLength, 0)
			})
		})

		Reset(func() {
			i18n.Close()
		})
	})
}
