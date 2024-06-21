package api

type ApiOption interface {
	apply(*Api)
}

type funcApiOption struct {
	f func(*Api)
}

func newFuncApiOption(f func(*Api)) *funcApiOption {
	return &funcApiOption{f}
}

func (o *funcApiOption) apply(api *Api) {
	o.f(api)
}
