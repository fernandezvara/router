package router

type Params struct {
	params map[string]string
	method string
	path   string
}

func newParams(params map[string]string, method, path string) *Params {
	return &Params{
		params: params,
		method: method,
		path:   path,
	}
}

func (s *Params) Param(key string) string {

	if _, ok := s.params[key]; ok {
		return s.params[key]
	}

	return ""

}

func (s *Params) Method() string {
	return s.method
}

func (s *Params) Path() string {
	return s.path
}
