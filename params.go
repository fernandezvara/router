package router

type Params struct {
	params map[string]string
	path   string
}

func newParams(params map[string]string, path string) *Params {
	return &Params{
		params: params,
		path:   path,
	}
}

func (s *Params) Param(key string) string {

	if _, ok := s.params[key]; ok {
		return s.params[key]
	}

	return ""

}

func (s *Params) Path() string {
	return s.path
}
