package test

import (
	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v0"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSimple(t *testing.T) {
	defer gock.Disable()
	gock.New("http://foo.com").
		Get("/bar").
		Reply(200).
		JSON(map[string]string{"foo": "bar"})

	res, err := http.Get("http://foo.com/bar")
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 200)

	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body)[:13], `{"foo":"bar"}`)
}
