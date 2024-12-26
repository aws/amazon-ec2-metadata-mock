package server

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var _ http.Handler = (*swapper)(nil)

// swapper is an `http.Handler` which allows the user to swap out routers while the server is running.
// This enables us to reset routes as the program is running, allowing for functionality like updating the config.
type swapper struct {
	mu     sync.Mutex
	router *mux.Router
}

func NewSwapper() *swapper {
	return &swapper{
		mu:     sync.Mutex{},
		router: mux.NewRouter(),
	}
}

func (rs *swapper) Reset() {
	rs.mu.Lock()
	rs.router = mux.NewRouter()
	rs.mu.Unlock()
}

func (rs *swapper) Walk(f func(route *mux.Route, r *mux.Router, ancestors []*mux.Route) error) {
	rs.mu.Lock()
	rs.router.Walk(f)
	rs.mu.Unlock()
}

func (rs *swapper) HandleFuncPrefix(pattern string, requestHandler HandlerType) {
	rs.mu.Lock()
	rs.router.PathPrefix(pattern).HandlerFunc(requestHandler)
	rs.mu.Unlock()
}

func (rs *swapper) HandleFunc(pattern string, requestHandler HandlerType) {
	rs.mu.Lock()
	rs.router.HandleFunc(pattern, requestHandler)
	rs.mu.Unlock()
}

func (rs *swapper) Swap(newRouter *mux.Router) {
	rs.mu.Lock()
	rs.router = newRouter
	rs.mu.Unlock()
}

func (rs *swapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rs.mu.Lock()
	router := rs.router
	rs.mu.Unlock()
	router.ServeHTTP(w, r)
}
