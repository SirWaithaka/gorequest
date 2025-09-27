package gorequest

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type FakeTemporaryError struct {
	error
	temporary bool
}

func (e FakeTemporaryError) Temporary() bool {
	return e.temporary
}

type MockHooks struct {
	str string
}

func (hooks *MockHooks) validate(*Request) {
	hooks.str = hooks.str + "validate:"
}

func (hooks *MockHooks) build(*Request) {
	hooks.str = hooks.str + "build:"
}

func (hooks *MockHooks) send(*Request) {
	hooks.str = hooks.str + "send:"
}

func (hooks *MockHooks) unmarshal(*Request) {
	hooks.str = hooks.str + "unmarshal:"
}

func (hooks *MockHooks) retry(*Request) {
	hooks.str = hooks.str + "retry:"
}

func (hooks *MockHooks) complete(*Request) {
	hooks.str = hooks.str + "complete:"
}

func TestRequest_New(t *testing.T) {

	t.Run("test retryer is not nil", func(t *testing.T) {
		req := New(Config{}, Operation{}, Hooks{}, nil, nil, nil)
		assert.NotNil(t, req.Retryer)
	})

	t.Run("test http request url", func(t *testing.T) {
		tcs := map[string]struct {
			Endpoint      string
			Path          string
			ExpectedPath  string
			ExpectedQuery string
		}{
			"no http Path": {
				Endpoint:      "https://example.com",
				Path:          "/",
				ExpectedPath:  "/",
				ExpectedQuery: "",
			},
			"with path in endpoint": {
				Endpoint:      "https://example.com/foo",
				Path:          "",
				ExpectedPath:  "/foo",
				ExpectedQuery: "",
			},
			"with query in path": {
				Endpoint:      "https://example.com",
				Path:          "/foo?bar=baz",
				ExpectedPath:  "/foo",
				ExpectedQuery: "bar=baz",
			},
			"with path in endpoint and query": {
				Endpoint:      "https://example.com/foo?bar=baz",
				Path:          "/qux",
				ExpectedPath:  "/foo/qux",
				ExpectedQuery: "",
			},
			"with query in path and endpoint": {
				Endpoint:      "https://example.com/?bar=baz",
				Path:          "/?foo=qux",
				ExpectedPath:  "/",
				ExpectedQuery: "foo=qux",
			},
		}

		for name, tc := range tcs {
			t.Run(name, func(t *testing.T) {
				op := Operation{Name: "FooBar", Path: tc.Path}
				req := New(Config{Endpoint: tc.Endpoint}, op, Hooks{}, nil, nil, nil)
				// assert results to expected values
				assert.Equal(t, tc.ExpectedPath, req.Request.URL.Path)
				assert.Equal(t, tc.ExpectedQuery, req.Request.URL.RawQuery)
				assert.Equal(t, http.MethodPost, req.Request.Method)
			})
		}
	})

	t.Run("test that it duplicates the hooks", func(t *testing.T) {
		// create a hook list with 4 hooks
		list := HookList{}
		list.PushBackHook(Hook{Name: "Foo", Fn: func(r *Request) {}})
		list.PushBackHook(Hook{Name: "Bar", Fn: func(r *Request) {}})
		list.PushBackHook(Hook{Name: "Baz", Fn: func(r *Request) {}})
		list.PushBackHook(Hook{Name: "Qux", Fn: func(r *Request) {}})
		// set hooks
		original := Hooks{}
		original.Validate = list.copy()
		original.Build = list.copy()
		original.Send = list.copy()
		original.Unmarshal = list.copy()
		original.Complete = list.copy()

		// make a test copy of the original hooks
		hooks := original.Copy()

		req := New(Config{}, Operation{}, hooks, nil, nil, nil)
		// assert the request hooks have not changed
		assert.Equal(t, hooks.Validate.Len(), req.Hooks.Validate.Len())
		assert.Equal(t, hooks.Build.Len(), req.Hooks.Build.Len())
		assert.Equal(t, hooks.Send.Len(), req.Hooks.Send.Len())
		assert.Equal(t, hooks.Unmarshal.Len(), req.Hooks.Unmarshal.Len())
		assert.Equal(t, hooks.Complete.Len(), req.Hooks.Complete.Len())

		// remove an item from each hook list in the test copy
		hooks.Validate.Remove("Foo")
		hooks.Validate.Remove("Bar")
		hooks.Build.Remove("Bar")
		hooks.Send.Remove("Baz")
		hooks.Unmarshal.Remove("Qux")
		//hooks.Complete.Remove("Qux")

		err := req.Send()
		assert.NoError(t, err)

		// assert the number of items in request hooks equals the original
		assert.Equal(t, original.Validate.Len(), req.Hooks.Validate.Len())
		assert.Equal(t, original.Build.Len(), req.Hooks.Build.Len())
		assert.Equal(t, original.Send.Len(), req.Hooks.Send.Len())
		assert.Equal(t, original.Unmarshal.Len(), req.Hooks.Unmarshal.Len())
		assert.Equal(t, original.Complete.Len(), req.Hooks.Complete.Len())

		// make a second request
		req = New(Config{}, Operation{}, hooks, nil, nil, nil)

		err = req.Send()
		assert.NoError(t, err)
		// assert the request hooks have not changed
		assert.Equal(t, hooks.Validate.Len(), req.Hooks.Validate.Len())
		assert.Equal(t, hooks.Build.Len(), req.Hooks.Build.Len())
		assert.Equal(t, hooks.Send.Len(), req.Hooks.Send.Len())
		assert.Equal(t, hooks.Unmarshal.Len(), req.Hooks.Unmarshal.Len())
		assert.Equal(t, hooks.Complete.Len(), req.Hooks.Complete.Len())

	})
}

func TestRequest_Send(t *testing.T) {

	t.Run("test that calling order of hooks is correct", func(t *testing.T) {

		// test that retry hooks are not called if no error occurs at Send hooks
		t.Run("test order when no error occurs at send hooks", func(t *testing.T) {
			mockHooks := MockHooks{}

			hooks := Hooks{
				Validate:  HookList{list: []Hook{{Fn: mockHooks.validate}}},
				Build:     HookList{list: []Hook{{Fn: mockHooks.build}}},
				Send:      HookList{list: []Hook{{Fn: mockHooks.send}}},
				Unmarshal: HookList{list: []Hook{{Fn: mockHooks.unmarshal}}},
				Retry:     HookList{list: []Hook{{Fn: mockHooks.retry}}},
				Complete:  HookList{list: []Hook{{Fn: mockHooks.complete}}},
			}

			req := New(Config{}, Operation{}, hooks.Copy(), nil, nil, nil)

			err := req.Send()
			assert.Nil(t, err)

			expected := "validate:build:send:unmarshal:complete:"
			assert.Equal(t, expected, mockHooks.str)
		})

		// test that retry hooks are called if an error occurs at Send hooks
		t.Run("test order when error occurs at send hooks", func(t *testing.T) {
			mockHooks := MockHooks{}

			hooks := Hooks{
				Validate:  HookList{list: []Hook{{Fn: mockHooks.validate}}},
				Build:     HookList{list: []Hook{{Fn: mockHooks.build}}},
				Send:      HookList{list: []Hook{{Fn: mockHooks.send}}},
				Unmarshal: HookList{list: []Hook{{Fn: mockHooks.unmarshal}}},
				Retry:     HookList{list: []Hook{{Fn: mockHooks.retry}}},
				Complete:  HookList{list: []Hook{{Fn: mockHooks.complete}}},
			}

			// mock an error at Send hooks
			hooks.Send.PushBack(func(r *Request) {
				// create a temporary error
				tempErr := FakeTemporaryError{error: errors.New("fake error"), temporary: true}
				r.Error = tempErr
			})
			req := New(Config{}, Operation{}, hooks, nil, nil, nil)

			err := req.Send()
			assert.NotNil(t, err)

			expected := "validate:build:send:complete:"
			assert.Equal(t, expected, mockHooks.str)
		})

	})

	t.Run("test that retryable requests are retried", func(t *testing.T) {
		hooks := Hooks{}

		cfg := RetryConfig{
			MaxRetries:     1,
			InitialDelay:   100 * time.Millisecond,
			Jitter:         0.1,
			MaxElapsedTime: 1 * time.Second,
		}

		// mock an error using send hook
		hooks.Send.PushBack(func(r *Request) {
			// create a temporary error
			tempErr := FakeTemporaryError{error: errors.New("fake error"), temporary: true}
			r.Error = tempErr
		})
		hooks.Retry.PushBack(func(r *Request) {
			r.RetryConfig.RetryCount++
		})

		// create an instance of retryer
		ret := retryer{}
		req := New(Config{}, Operation{}, hooks, ret, nil, nil)
		req.WithRetryConfig(cfg)

		err := req.Send()
		assert.NotNil(t, err)

		// confirm that request was retried
		assert.Equal(t, cfg.MaxRetries, req.RetryConfig.RetryCount)
	})
}
