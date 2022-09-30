package appfinger

import "testing"

func TestT(t *testing.T) {
	convBanner(struct {
		Host   string
		AAA    string
		Header string
	}{
		Host:   "asdfasdfasd",
		AAA:    "asdfasdfasd",
		Header: "11111",
	})
}
