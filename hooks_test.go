package gorequest_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SirWaithaka/gorequest"
)

func TestHookList_Run(t *testing.T) {
	r := &gorequest.Request{}
	h := gorequest.HookList{}

	val := ""
	h.PushBack(func(r *gorequest.Request) {
		val += "a"
		r.Params = val
	})
	h.Run(r)

	// assert
	assert.Equal(t, "a", val)
	assert.Equal(t, "a", r.Params)
}

func TestHooksList_Remove(t *testing.T) {
	hooks := gorequest.HookList{}
	hook := gorequest.Hook{Name: "Foo", Fn: func(r *gorequest.Request) {}}
	hook2 := gorequest.Hook{Name: "Bar", Fn: func(r *gorequest.Request) {}}
	// add 4 hooks
	hooks.PushFrontHook(hook)
	hooks.PushFrontHook(hook2)
	hooks.PushFrontHook(hook)
	hooks.PushFront(func(r *gorequest.Request) {})

	// assert for 4 hooks
	assert.Equal(t, 4, hooks.Len())

	// remove hook
	hooks.RemoveHook(hook)
	assert.Equal(t, 2, hooks.Len())

}
