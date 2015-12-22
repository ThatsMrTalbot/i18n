package i18n

import (
	"net/http"
	"testing"

	"golang.org/x/text/language"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShard(t *testing.T) {
	t.Parallel()

	Convey("Given given many requests pointer", t, func() {
		reqs := make([]*http.Request, 100)
		for i := 0; i < 100; i++ {
			reqs[i], _ = http.NewRequest("GET", "http://example.com", nil)
		}

		Convey("When the requests are sharded", func() {
			var empty struct{}
			shards := make(map[string]struct{})

			for _, req := range reqs {
				str := shardGenerate(req)
				shards[str] = empty
			}

			Convey("Then there should be a range of shard values", func() {
				//Printf("Shard: 100 requests sharded into %d values", len(shards))
				So(len(shards), ShouldBeGreaterThan, 50)
			})
		})
	})

	Convey("Given given a request pointer", t, func() {
		req, _ := http.NewRequest("GET", "http://example.com", nil)

		Convey("When the requests is stored", func() {
			m := newRequestLanguageMap()
			m.Add(req, language.English)

			Convey("Then the request language should be stored", func() {
				lang, ok := m.Get(req)
				So(lang.String(), ShouldEqual, language.English.String())
				So(ok, ShouldBeTrue)
			})
		})
	})

	Convey("Given given a request map", t, func() {
		req, _ := http.NewRequest("GET", "http://example.com", nil)

		Convey("When the requests is deleted", func() {
			m := newRequestLanguageMap()
			m.Add(req, language.English)
			m.Delete(req)

			Convey("Then the request language should deleted", func() {
				lang, ok := m.Get(req)
				So(lang.String(), ShouldEqual, language.Und.String())
				So(ok, ShouldBeFalse)
			})
		})
	})
}
