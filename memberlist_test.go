package ringman

import (
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Nitro/memberlist"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewMemberlistRing(t *testing.T) {
	// This must be started out here so we don't re-use the port and
	// start the list up each time
	mlConfig := memberlist.DefaultLANConfig()
	mlConfig.BindPort = 35001

	Convey("NewMemberlistRing()", t, func() {
		mlistRing, err := NewMemberlistRing(mlConfig, []string{}, "8000", "default")
		So(err, ShouldBeNil)

		ourAddr := mlistRing.Memberlist.LocalNode().Addr.String()
		Convey("returns a properly configured MemberlistRing", func() {
			So(mlistRing.Memberlist, ShouldNotBeNil)
			So(mlistRing.manager, ShouldNotBeNil)

			node, err := mlistRing.manager.GetNode("beowulf")
			So(err, ShouldBeNil)
			So(node, ShouldEqual, ourAddr+":8000")

			So(len(mlistRing.Memberlist.Members()), ShouldEqual, 1)
		})
	})
}

func Test_MemberListRingShutdown(t *testing.T) {
	Convey("NewMemberlistRing()", t, func() {
		mlistRing, err := NewDefaultMemberlistRing([]string{}, "8000")
		So(err, ShouldBeNil)

		mlistRing.Shutdown()

		So(mlistRing.manager.Ping(), ShouldBeFalse)
		So(mlistRing.manager.cmdChan, ShouldBeNil)
	})
}

func Test_HttpGetNodeHandler(t *testing.T) {
	Convey("HttpGetNodeHandler()", t, func() {
		// Don't initialize it because we don't need it to be
		mlistRing := &MemberlistRing{}

		req := httptest.NewRequest("GET", "/services/boccacio.json", nil)
		recorder := httptest.NewRecorder()

		Convey("returns a 404 when no key is provided", func() {
			mlistRing.HttpGetNodeHandler(recorder, req)

			So(recorder.Result().StatusCode, ShouldEqual, 404)
		})

		Convey("returns a node when a key is provided", func() {
			form := url.Values{}
			form.Set("key", "bocaccio")
			req.Form = form

			mlistRing.HttpGetNodeHandler(recorder, req)

			bodyBytes, _ := ioutil.ReadAll(recorder.Result().Body)
			body := string(bodyBytes)

			So(recorder.Result().StatusCode, ShouldEqual, 200)
			So(body, ShouldContainSubstring, `"Key": "bocaccio"`)
		})
	})
}
