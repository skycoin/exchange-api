package c2cx

import (
	"net/url"
	"testing"
)

func Test_sign(t *testing.T) {
	var params = url.Values{}
	params.Add("apiKey", "C821DB84-6FBD-11E4-A9E3-C86000D26D7C")
	want := "BC0DE7EBA50C730BDFC575FE2CD54082"
	expected := sign("12D857DE-7A92-F555-10AC-7566A0D84D1B", params)
	if want != expected {
		t.Fatalf("Incorrect sign!\nwant %s, expected %s", want, expected)
	}
}
