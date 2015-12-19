package i18n

import (
	"testing"

	"golang.org/x/text/language"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInMemoryStorage(t *testing.T) {
	t.Parallel()

	Convey("Given an empty memory store", t, func() {
		storage := NewInMemoryStorage()
		Convey("When an item is added to the memory store", func() {

			expected := &Translation{
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

			expected := &Translation{
				Lang:  language.English,
				Key:   "SomeKey",
				Value: "SomeValue",
			}

			replacement := &Translation{
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
		storage := NewInMemoryStorage()
		err := storage.Store(&Translation{
			Lang:  language.English,
			Key:   "SomeKey",
			Value: "SomeValue",
		})
		So(err, ShouldBeNil)

		Convey("When an item is deleted from the storage", func() {
			storage.Delete(&Translation{
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
