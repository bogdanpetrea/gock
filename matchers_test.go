package gock

import (
	"github.com/nbio/st"
	"net/http"
	"net/url"
	"testing"
)

func TestMatchMethod(t *testing.T) {
	cases := []struct {
		value   string
		method  string
		matches bool
	}{
		{"GET", "GET", true},
		{"POST", "POST", true},
		{"", "POST", true},
		{"POST", "GET", false},
		{"PUT", "GET", false},
	}

	for _, test := range cases {
		req := &http.Request{Method: test.method}
		ereq := &Request{Method: test.value}
		matches, err := MatchMethod(req, ereq)
		st.Expect(t, err, nil)
		st.Expect(t, matches, test.matches)
	}
}

func TestMatchHost(t *testing.T) {
	cases := []struct {
		value   string
		url     string
		matches bool
	}{
		{"foo.com", "foo.com", true},
		{"foo.net", "foo.com", false},
		{"foo", "foo.com", true},
		{"(.*).com", "foo.com", true},
		{"127.0.0.1", "127.0.0.1", true},
		{"127.0.0.2", "127.0.0.1", false},
		{"127.0.0.*", "127.0.0.1", true},
		{"127.0.0.[0-9]", "127.0.0.7", true},
	}

	for _, test := range cases {
		req := &http.Request{URL: &url.URL{Host: test.url}}
		ereq := &Request{URLStruct: &url.URL{Host: test.value}}
		matches, err := MatchHost(req, ereq)
		st.Expect(t, err, nil)
		st.Expect(t, matches, test.matches)
	}
}

func TestMatchPath(t *testing.T) {
	cases := []struct {
		value   string
		path    string
		matches bool
	}{
		{"/foo", "/foo", true},
		{"/foo", "/foo/bar", true},
		{"bar", "/foo/bar", true},
		{"foo", "/foo/bar", true},
		{"bar$", "/foo/bar", true},
		{"/foo/*", "/foo/bar", true},
		{"/foo/[a-z]+", "/foo/bar", true},
		{"/foo/baz", "/foo/bar", false},
		{"/foo/baz", "/foo/bar", false},
	}

	for _, test := range cases {
		u, _ := url.Parse("http://foo.com" + test.path)
		mu, _ := url.Parse("http://foo.com" + test.value)
		req := &http.Request{URL: u}
		ereq := &Request{URLStruct: mu}
		matches, err := MatchPath(req, ereq)
		st.Expect(t, err, nil)
		st.Expect(t, matches, test.matches)
	}
}

func TestMatchHeaders(t *testing.T) {
	cases := []struct {
		values  http.Header
		headers http.Header
		matches bool
	}{
		{http.Header{"foo": []string{"bar"}}, http.Header{"foo": []string{"bar"}}, true},
		{http.Header{"foo": []string{"bar"}}, http.Header{"foo": []string{"barbar"}}, true},
		{http.Header{"bar": []string{"bar"}}, http.Header{"foo": []string{"bar"}}, false},
		{http.Header{"foofoo": []string{"bar"}}, http.Header{"foo": []string{"bar"}}, false},
		{http.Header{"foo": []string{"bar(.*)"}}, http.Header{"foo": []string{"barbar"}}, true},
		{http.Header{"foo": []string{"b(.*)"}}, http.Header{"foo": []string{"barbar"}}, true},
		{http.Header{"foo": []string{"^bar$"}}, http.Header{"foo": []string{"bar"}}, true},
		{http.Header{"foo": []string{"^bar$"}}, http.Header{"foo": []string{"barbar"}}, false},
	}

	for _, test := range cases {
		req := &http.Request{Header: test.headers}
		ereq := &Request{Header: test.values}
		matches, err := MatchHeaders(req, ereq)
		st.Expect(t, err, nil)
		st.Expect(t, matches, test.matches)
	}
}

func TestMatchBody(t *testing.T) {
	cases := []struct {
		value   string
		body    string
		matches bool
	}{
		{"foo bar", "foo bar\n", true},
		{"foo", "foo bar\n", true},
		{"f[o]+", "foo\n", true},
		{`"foo"`, `{"foo":"bar"}\n`, true},
		{`{"foo":"bar"}`, `{"foo":"bar"}\n`, true},
		{`{"foo":"foo"}`, `{"foo":"bar"}\n`, false},
	}

	for _, test := range cases {
		req := &http.Request{Body: createReadCloser([]byte(test.body))}
		ereq := &Request{BodyBuffer: []byte(test.value)}
		matches, err := MatchBody(req, ereq)
		st.Expect(t, err, nil)
		st.Expect(t, matches, test.matches)
	}
}
