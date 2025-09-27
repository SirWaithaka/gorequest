package corehooks_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"

	"github.com/SirWaithaka/gorequest"
	"github.com/SirWaithaka/gorequest/corehooks"
)

func TestAddScheme(t *testing.T) {
	tcs := map[string]struct {
		Endpoint   string
		DisableSSL bool
		Expected   string
	}{
		"with no scheme": {
			Endpoint:   "example.com",
			DisableSSL: false,
			Expected:   "https://example.com",
		},
		"disable ssl": {
			Endpoint:   "example.com",
			DisableSSL: true,
			Expected:   "http://example.com",
		},
		"with ssl scheme": {
			Endpoint:   "https://example.com",
			DisableSSL: false,
			Expected:   "https://example.com",
		},
		"with no ssl scheme": {
			Endpoint:   "http://example.com",
			DisableSSL: false,
			Expected:   "http://example.com",
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			endpoint := corehooks.AddScheme(tc.Endpoint, tc.DisableSSL)
			assert.Equal(t, tc.Expected, endpoint)
		})
	}
}

type testSendHandlerTransport struct{ timeout time.Duration }

func (t *testSendHandlerTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if t.timeout > 0 {
		time.Sleep(t.timeout)
	}
	return nil, errors.New("mock error")
}

func TestSendHook(t *testing.T) {

	t.Run("test redirect", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/redirect":
				u := *r.URL
				u.Path = "/home"
				w.Header().Set("Location", u.String())
				w.WriteHeader(http.StatusTemporaryRedirect)

			case "/home":
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("ok"))
			}
		}))
		defer server.Close()

		tcs := map[string]struct {
			Redirect       bool
			ExpectedStatus int
		}{
			"redirect": {
				Redirect:       true,
				ExpectedStatus: http.StatusOK,
			},
			"no redirect": {
				Redirect:       false,
				ExpectedStatus: http.StatusTemporaryRedirect,
			},
		}

		for name, tc := range tcs {
			t.Run(name, func(t *testing.T) {
				cfg := gorequest.Config{Endpoint: server.URL, DisableSSL: true, HTTPClient: http.DefaultClient}
				op := gorequest.Operation{Name: "FooBar", Path: "/redirect"}

				hooks := gorequest.Hooks{}
				hooks.Send.PushBackHook(corehooks.SendHook)

				cfg.DisableFollowRedirects = !tc.Redirect

				req := gorequest.New(cfg, op, hooks, nil, nil, nil)
				err := req.Send()
				assert.Nil(t, err)

				// check response status
				assert.Equal(t, tc.ExpectedStatus, req.Response.StatusCode)

			})
		}
	})

	t.Run("test handle send error", func(t *testing.T) {

		t.Run("transport error", func(t *testing.T) {
			client := &http.Client{Transport: &testSendHandlerTransport{}}
			op := gorequest.Operation{Name: "Operation"}

			hooks := gorequest.Hooks{}
			hooks.Send.PushBackHook(corehooks.SendHook)
			req := gorequest.New(gorequest.Config{HTTPClient: client}, op, hooks, nil, nil, nil)

			err := req.Send()
			assert.NotNil(t, err)
			assert.NotNil(t, req.Response)
		})

		t.Run("url.Error timeout", func(t *testing.T) {
			client := &http.Client{
				Timeout:   100 * time.Millisecond,
				Transport: &testSendHandlerTransport{timeout: 500 * time.Millisecond},
			}
			op := gorequest.Operation{Name: "Operation"}

			hooks := gorequest.Hooks{}
			hooks.Send.PushBackHook(corehooks.SendHook)
			req := gorequest.New(gorequest.Config{HTTPClient: client}, op, hooks, nil, nil, nil)

			err := req.Send()
			assert.NotNil(t, err)
			assert.NotNil(t, req.Response)
			assert.Equal(t, 0, req.Response.StatusCode)
		})
	})
}

func TestSetRequestID(t *testing.T) {
	rid := xid.New().String()

	generator := func() string {
		return rid
	}

	// build request hooks
	hooks := gorequest.Hooks{}
	hooks.Build.PushFrontHook(corehooks.SetRequestID(generator))
	hooks.Complete.PushFront(func(r *gorequest.Request) {
		assert.Equal(t, rid, r.Config.RequestID)
	})

	req := gorequest.New(gorequest.Config{}, gorequest.Operation{}, hooks, nil, nil, nil)
	err := req.Send()
	assert.Nil(t, err)
}
