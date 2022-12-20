package router

import (
	"strings"
)

const (
	pathDelimiter string = "/"
)

type Tree struct {
	leaf     *Leaf
	method   string
	notFound Handler
}

// NewTree creates a new trie tree.
func NewTree(method string, notFound Handler) *Tree {
	return &Tree{
		leaf:     newLeaf(""),
		method:   method,
		notFound: notFound,
	}
}

func (t *Tree) valid(path string) (err error) {

	var ok bool
	if ok = validRegExp.MatchString(path); !ok {
		err = ErrInvalidRoute
	}

	return

}

func (t *Tree) Insert(path string, handler Handler) error {

	var (
		leaf     *Leaf
		ok       bool
		parts    []string
		index    int
		thisPath string
		err      error
	)

	if err = t.valid(path); err != nil {
		return err
	}

	leaf = t.leaf
	parts = strings.Split(path, pathDelimiter)

	for index, thisPath = range parts {

		var subleaf *Leaf
		subleaf, ok = leaf.Static[thisPath]

		if ok {
			// node exists
			leaf = subleaf
			continue
		} else {

			// is dynamic?
			if strings.HasPrefix(thisPath, ":") {
				// dynamic leaf already exists?
				if subleaf, ok = leaf.Dynamic[thisPath]; ok {
					leaf = subleaf
					continue
				}

				// create new dynamic node
				var newLeaf = newLeaf(thisPath)
				newLeaf.Param = thisPath[1:]
				leaf.Dynamic[thisPath] = newLeaf
				leaf = newLeaf

			} else {
				// create new static node
				var newLeaf = newLeaf(thisPath)
				leaf.Static[thisPath] = newLeaf
				leaf = newLeaf

			}

		}

		// last loop.
		if index == len(parts)-1 {
			leaf.Path = path
			leaf.Handler = handler
			break
		}

	}

	return nil

}

func (t *Tree) Execute(path string) error {

	var (
		p   *Params
		h   Handler
		err error
	)

	if p, h, err = t.search(path); err != nil {
		return err
	}

	return h(p)

}

func (t *Tree) search(path string) (p *Params, h Handler, err error) {

	var (
		leaf   *Leaf = t.leaf
		parts  []string
		index  int
		params map[string]string = make(map[string]string)
		// retParams map[string]string
	)

	parts = strings.Split(path, pathDelimiter)

	params, h, err = t.subsearch(path, leaf, parts, index, params)

	p = newParams(params, t.method, path)

	if err == ErrNotFound {
		err = t.notFound(p)
	}

	return

}

func (t *Tree) subsearch(path string, leaf *Leaf, parts []string, index int, params map[string]string) (retParams map[string]string, h Handler, err error) {

	var (
		subleaf *Leaf
		ok      bool
	)

	retParams = params

	if subleaf, ok = leaf.Static[parts[index]]; ok {
		if len(parts)-1 == index {
			h = subleaf.Handler
			return
		}

		return t.subsearch(path, subleaf, parts, index+1, params)
	}

	if retParams, h, err = t.subDynSearch(path, leaf, parts, index, params); err == nil {
		return
	}

	err = ErrNotFound
	return

}

func (t *Tree) subDynSearch(path string, leaf *Leaf, parts []string, index int, params map[string]string) (retParams map[string]string, h Handler, err error) {

	retParams = params

	for _, subleaf := range leaf.Dynamic {

		// set the param
		retParams[subleaf.Param] = parts[index]

		// is last subleaf dynamic?
		if len(parts)-1 == index {
			h = subleaf.Handler
			return

		}

		// else, continue searching using this subpart
		if retParams, h, err = t.subsearch(path, subleaf, parts, index+1, params); err == nil {
			// return only if found
			return
		}

	}

	err = ErrNotFound
	return

}
