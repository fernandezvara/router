package router

type Leaf struct {
	Path    string
	Handler Handler
	Static  map[string]*Leaf
	Dynamic map[string]*Leaf
	Param   string // parameter name
}

func newLeaf(path string) *Leaf {
	return &Leaf{
		Path:    path,
		Static:  make(map[string]*Leaf),
		Dynamic: make(map[string]*Leaf),
	}
}
