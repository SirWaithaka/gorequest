package gorequest_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SirWaithaka/gorequest"
)

func TestHookList_Copy(t *testing.T) {

	// create a hook list with 4 hooks
	list := gorequest.HookList{}
	list.PushBackHook(gorequest.Hook{Name: "Foo", Fn: func(r *gorequest.Request) {}})
	list.PushBackHook(gorequest.Hook{Name: "Bar", Fn: func(r *gorequest.Request) {}})
	list.PushBackHook(gorequest.Hook{Name: "Baz", Fn: func(r *gorequest.Request) {}})
	list.PushBackHook(gorequest.Hook{Name: "Qux", Fn: func(r *gorequest.Request) {}})
	// set hooks
	hooks := gorequest.Hooks{}
	hooks.Validate = list
	hooks.Build = list
	hooks.Send = list
	hooks.Unmarshal = list
	hooks.Complete = list
	// create a copy
	copied := hooks.Copy()
	// assert that the number of items in copy and original are the same
	assert.Equal(t, hooks.Validate.Len(), copied.Validate.Len())
	assert.Equal(t, hooks.Build.Len(), copied.Build.Len())
	assert.Equal(t, hooks.Send.Len(), copied.Send.Len())
	assert.Equal(t, hooks.Unmarshal.Len(), copied.Unmarshal.Len())
	assert.Equal(t, hooks.Complete.Len(), copied.Complete.Len())

	// delete an item from the original
	hooks.Validate.Remove("Bar")
	hooks.Build.Remove("Baz")
	hooks.Send.Remove("Qux")
	hooks.Unmarshal.Remove("Foo")
	hooks.Complete.Remove("Qux")
	// assert that the copy has the removed item
	assert.Equal(t, 4, copied.Validate.Len())
	assert.Equal(t, 4, copied.Build.Len())
	assert.Equal(t, 4, copied.Send.Len())
	assert.Equal(t, 4, copied.Unmarshal.Len())
	assert.Equal(t, 4, copied.Complete.Len())

}

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
