package main

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {
	setupServers()
	defer teardownServers()

	client := setupClient()

	Convey("Given a running hustle ws server", t, func() {
		Convey("When publishing", func() {
			done := make(chan bool)
			go func() {
				err := client.Publish("test", "test", "test")
				So(err, ShouldNotEqual, nil)
				done <- true
			}()

			select {
			case <-done:
				return
			case <-time.After(5 * time.Second):
				t.Fail()
			}
		})
	})
}
