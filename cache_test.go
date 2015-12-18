package i18n

import (
	"testing"

	"golang.org/x/text/language"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCache(t *testing.T) {
	t.Parallel()

	Convey("Given an empty cache", t, func() {
		cache := new(Cache)
		Convey("When an item is added to the cache", func() {
			expected := &Translation{
				Lang: language.English,
				Key:    "SomeKey",
				Value:  "SomeValue",
			}
			cache.Add(expected)

			Convey("Then it should be accessable", func() {
				result := cache.Get(language.English, "SomeKey")
				So(result, ShouldNotBeNil)
				So(result, ShouldResemble, expected)
			})
		})
	})

	Convey("Given a populated cache", t, func() {
		cache := new(Cache)
		cache.Add(&Translation{
			Lang: language.English,
			Key:    "SomeKey",
			Value:  "SomeValue",
		})

		Convey("When an item is deleted from the cache", func() {
			cache.Delete(&Translation{
				Lang: language.English,
				Key:    "SomeKey",
				Value:  "SomeValue",
			})

			Convey("The item should not exists in the cache", func() {
				result := cache.Get(language.English, "SomeKey")
				So(result, ShouldBeNil)
			})
		})

		Convey("When the cache is cleared", func() {
			cache.Clear()

			Convey("The item should not exists in the cache", func() {
				result := cache.Get(language.English, "SomeKey")
				So(result, ShouldBeNil)
			})
		})
	})
}
