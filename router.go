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
	methods  map[string]*Tree
	notFound Handler
}

// New returns a new router
func New(notFound Handler) *Router {

	var r Router = Router{
		methods:  make(map[string]*Tree),
		notFound: notFound,
	}

	if r.notFound == nil {
		r.notFound = defaultNotFound
	}

	return &r

}

// Method returns the tree with the handlers for this 'method'
func (r *Router) Method(method string) *Tree {

	if _, ok := r.methods[method]; !ok {
		r.methods[method] = NewTree(method, r.notFound)
	}

	return r.methods[method]

}

func defaultNotFound(_ *Params) error {

	return ErrNotFound

}
