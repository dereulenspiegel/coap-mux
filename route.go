package mux

import (
	"net"

	coap "github.com/dustin/go-coap"
)

type RouteMatch struct {
	Handler coap.Handler
	Vars    map[string]string
}

type Matcher interface {
	Match(*coap.Message, *net.UDPAddr) bool
}

type Route struct {
	name     string
	handler  coap.Handler
	matchers []Matcher
	regexp   *routeRegexp
}

// getRegexpGroup returns regexp definitions from this route.
func (r *Route) getRegexp() *routeRegexp {
	if r.regexp == nil {
		r.regexp = new(routeRegexp)
	}
	return r.regexp
}

func (r *Route) Name(name string) *Route {
	r.name = name
	return r
}

func (r *Route) Match(msg *coap.Message, addr *net.UDPAddr, match *RouteMatch) bool {
	for _, matcher := range r.matchers {
		if matched := matcher.Match(msg, addr); !matched {
			return false
		}
	}
	match.Handler = r.handler
	if match.Vars == nil {
		match.Vars = make(map[string]string)
	}
	// Set variables.
	if r.regexp != nil {
		r.regexp.setMatch(msg, match, r)
	}
	return true
}

func (r *Route) Handler(h coap.Handler) *Route {
	r.handler = h
	return r
}

// addMatcher adds a matcher to the route.
func (r *Route) addMatcher(m Matcher) *Route {
	r.matchers = append(r.matchers, m)
	return r
}

func (r *Route) Matches(matcher Matcher) *Route {
	r.addMatcher(matcher)
	return r
}

func (r *Route) Methods(methods ...coap.COAPCode) *Route {
	return r.addMatcher(methodMatcher(methods))
}

func (r *Route) Path(tpl string) *Route {
	r.addRegexpMatcher(tpl)
	return r
}
