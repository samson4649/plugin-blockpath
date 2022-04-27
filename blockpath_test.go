package plugin_blockpath

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		desc    string
		// regexps []string
		elements []Element
		expErr  bool
	}{
		{
			desc:    "should return no error",
			elements: []Element{
				Element{
					Regex: `^/foo/(.*)`,
					Response: http.StatusForbidden,
				},
			},
			expErr:  false,
		},
		{
			desc:    "should return an error",
			elements: []Element{
				Element{
					Regex: `*`,
					Response: http.StatusForbidden,
				},
			},
			expErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cfg := &Config{
				Elements: test.elements,
			}

			if _, err := New(context.Background(), nil, cfg, "name"); test.expErr && err == nil {
				t.Errorf("expected error on bad regexp format")
			}
		})
	}
}

func TestServeHTTP(t *testing.T) {
	tests := []struct {
		desc          string
		// regexps       []string
		elements      []Element
		reqPath       string
		expNextCall   bool
		expStatusCode int
	}{
		{
			desc:          "should return forbidden status",
			elements: []Element{
				Element{
					Regex: `/test`,
					Response: http.StatusForbidden,
				},
			},
			// regexps:       []string{"/test"},
			reqPath:       "/test",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:          "should return forbidden status",
			elements: []Element{
				Element{
					Regex: `/test`,
					Response: http.StatusForbidden,
				},
				Element{
					Regex: `/toto`,
					Response: http.StatusForbidden,
				},
			},
			// regexps:       []string{"/test", "/toto"},
			reqPath:       "/toto",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:          "should return ok status",
			elements: []Element{
				Element{
					Regex: `/test`,
					Response: http.StatusForbidden,
				},
				Element{
					Regex: `/toto`,
					Response: http.StatusForbidden,
				},
			},
			// regexps:       []string{"/test", "/toto"},
			reqPath:       "/plop",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:          "should return not found status",
			elements: []Element{
				Element{
					Regex: `/test`,
					Response: http.StatusNotFound,
				},
				Element{
					Regex: `/toto`,
					Response: http.StatusNotFound,
				},
			},
			// regexps:       []string{"/test", "/toto"},
			reqPath:       "/test",
			expNextCall:   false,
			expStatusCode: http.StatusNotFound,
		},
		{
			desc:          "should return ok status",
			reqPath:       "/test",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:          "should return forbidden status",
			elements: []Element{
				Element{
					Regex: `^/bar(.*)`,
					Response: http.StatusForbidden,
				},
			},
			// regexps:       []string{`^/bar(.*)`},
			reqPath:       "/bar/foo",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:          "should return forbidden status",
			elements: []Element{
				Element{
					Regex: `^/bar(.*)`,
					Response: http.StatusForbidden,
				},
			},
			// regexps:       []string{`^/bar(.*)`},
			reqPath:       "/foo/bar",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cfg := &Config{
				Elements: test.elements,
			}

			nextCall := false
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				nextCall = true
			})

			handler, err := New(context.Background(), next, cfg, "blockpath")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("http://localhost%s", test.reqPath)
			req := httptest.NewRequest(http.MethodGet, url, nil)

			handler.ServeHTTP(recorder, req)

			if nextCall != test.expNextCall {
				t.Errorf("next handler should not be called")
			}

			if recorder.Result().StatusCode != test.expStatusCode {
				t.Errorf("got status code %d, want %d", recorder.Code, test.expStatusCode)
			}
		})
	}
}
