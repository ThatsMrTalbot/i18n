package server

import (
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
