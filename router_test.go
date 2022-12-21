package router

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {

	r := New(nil).Method("TEST")
	assert.NotNil(t, r)

	r.Insert("m/:a", func(_ *Params) error { print("m/:a"); return nil })
	assert.Equal(t, r.leaf.Path, "")

	_, _, e := r.Search("m/a")
	assert.NotNil(t, r)
	assert.Nil(t, e)

	for _, u := range []string{
		"items",
		"items/:id",
		"items/:id/subitems",
		"items/:id/subitems/:sid",
		"items/:id/sub",
		"items/:id/sub/:sid",
	} {
		assert.Nil(t, r.Insert(u, func(p *Params) error {
			print(u)
			fmt.Printf("p: %#+v\n", p)
			return nil
		}))
	}

	// invalid path
	assert.Error(t, r.Insert("this contains spaces so is not valid", func(_ *Params) error { return nil }))
	assert.Error(t, r.Insert("this,contains,commas,so,is,not,valid", func(_ *Params) error { return nil }))

	_, _, e = r.Search("items")
	assert.Nil(t, e)

	_, _, e = r.Search("items/1")
	assert.Nil(t, e)

	_, _, e = r.Search("items/1/subitems")
	assert.Nil(t, e)

	_, _, e = r.Search("items/1/subitems/11")
	assert.Nil(t, e)

	_, _, e = r.Search("items/1/sub")
	assert.Nil(t, e)

	_, _, e = r.Search("items/1/sub/11")
	assert.Nil(t, e)

	_, _, e = r.Search(("non-existent"))
	assert.Equal(t, ErrNotFound, e)

	// test handler

	err := r.Insert("hh/:a/:b", func(p *Params) error {

		assert.Equal(t, "aaa", p.Param("a"))
		assert.Equal(t, "bbbb", p.Param("b"))
		assert.Equal(t, "hh/aaa/bbbb", p.Path())
		assert.Equal(t, "", p.Param("non-existent"))

		return nil

	})
	assert.Nil(t, err)

	err = r.Execute("hh/aaa/bbbb")
	assert.Nil(t, err)

	err = r.Execute("nonexistent")
	assert.Equal(t, ErrNotFound, err)

}

func TestCustomNotFound(t *testing.T) {

	var errCustom = errors.New("my custom error")

	custom := func(p *Params) error {
		// params are accesible here to help find potential missing routes, log, debug...
		println(p.Method())
		println(p.Path())
		return errCustom
	}

	r := New(custom).Method("TEST")
	assert.NotNil(t, r)

	r.Insert("m/:a", func(_ *Params) error { print("m/:a"); return nil })
	assert.Equal(t, r.leaf.Path, "")

	params, _, e := r.Search("m/test")
	assert.NotNil(t, r)
	assert.Nil(t, e)
	assert.Equal(t, "test", params.Param("a"))
	assert.Equal(t, "TEST", params.Method())
	assert.Equal(t, "m/test", params.Path())

	params, _, e = r.Search("path/that/not/exists")
	assert.Equal(t, errCustom, e)
	assert.Equal(t, "TEST", params.Method())
	assert.Equal(t, "path/that/not/exists", params.Path())

	e = r.Execute("notexists")
	assert.Equal(t, errCustom, e)

}
