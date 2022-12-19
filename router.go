package router

import (
	"errors"
	"regexp"
)

var (
	// errors
	ErrNotFound     = errors.New("not found")
	ErrInvalidRoute = errors.New("invalid route")

	// validation
	validRegExp = regexp.MustCompile("^[a-zA-Z0-9:/]+$")
)

type Handler func(*Params) error

// Router is a request router for any kind of workload, not specific for HTTP servers
type Router struct {
	methods map[string]*Tree
}

// New returns a new router
func New() *Router {

	return &Router{
		methods: make(map[string]*Tree),
	}

}

// Method returns the tree with the handlers for this 'method'
func (r *Router) Method(method string) *Tree {

	if _, ok := r.methods[method]; !ok {
		r.methods[method] = NewTree()
	}

	return r.methods[method]

}
