package mux

import (
	"net"
	"sync"

	coap "github.com/dustin/go-coap"
)

var (
	mutex   sync.RWMutex
	msgVars = make(map[*coap.Message]map[string]string)
)

func SetVar(msg *coap.Message, name, data string) {
	mutex.Lock()
	r, exists := msgVars[msg]
	if !exists {
		r = make(map[string]string)
	}
	r[name] = data
	mutex.Unlock()
}

func SetVars(msg *coap.Message, vars map[string]string) {
	mutex.Lock()
	msgVars[msg] = vars
	mutex.Unlock()
}

func Var(msg *coap.Message, name string) string {
	mutex.RLock()
	var value = ""

	r, exists := msgVars[msg]
	if exists {
		value = r[name]
	}
	mutex.RUnlock()
	return value
}

func ClearVars(msg *coap.Message) {
	mutex.Lock()
	delete(msgVars, msg)
	mutex.Unlock()
}

type Router struct {
	NotFoundHandler coap.Handler
	routes          []*Route
}

func NewRouter() *Router {
	return &Router{routes: make([]*Route, 0, 50)}
}

func (r *Router) Match(msg *coap.Message, addr *net.UDPAddr, match *RouteMatch) bool {
	for _, route := range r.routes {
		if route.Match(msg, addr, match) {
			return true
		}
	}
	return false
}

func (r *Router) NewRoute() *Route {
	route := &Route{}
	r.routes = append(r.routes, route)
	return route
}

func (r *Router) Handle(path string, handler coap.Handler) *Route {
	return r.NewRoute().Path(path).Handler(handler)
}

func (r *Router) Path(tpl string) *Route {
	return r.NewRoute().Path(tpl)
}

func (r *Router) ServeCOAP(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	var match *RouteMatch
	var returnMessage *coap.Message
	if r.Match(m, a, match) {
		// TODO set vars
		SetVars(m, match.Vars)
		returnMessage = match.Handler.ServeCOAP(l, a, m)
		ClearVars(m)
	} else {
		returnMessage = r.NotFoundHandler.ServeCOAP(l, a, m)
	}
	return returnMessage
}
