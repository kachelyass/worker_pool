package httpm

import (
	"net/http"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

func (r *Router) Handle(route string, handler http.Handler) {
	r.mux.Handle(route, Middleware(route, handler))
}

func (r *Router) HandleFunc(route string, fn http.HandlerFunc) {
	r.Handle(route, fn)
}

func (r *Router) RawHandle(route string, handler http.Handler) {
	r.mux.Handle(route, handler)
}

func (r *Router) RawHandleFunc(route string, fn http.HandlerFunc) {
	r.mux.HandleFunc(route, fn)
}

func (r *Router) Handler() http.Handler {
	return r.mux
}
