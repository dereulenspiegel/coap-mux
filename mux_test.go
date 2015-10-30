package mux

import (
	"testing"

	"github.com/stretchr/testify/assert"

	coap "github.com/dustin/go-coap"
)

func TestBasicRouting(t *testing.T) {
	assert := assert.New(t)

	router := NewRouter()
	router.Handle("/a/{id1}/{id2}/b/{id3}", nil).Methods(coap.GET)

	validMsg := &coap.Message{}
	validMsg.SetPathString("/a/1234/abcd/b/23")
	validMsg.Code = coap.GET

	match := &RouteMatch{}
	assert.True(router.Match(validMsg, nil, match))

	invalidMsg := &coap.Message{}
	invalidMsg.SetPathString("/a/1234/abcd/b/23")
	invalidMsg.Code = coap.PUT

	match = &RouteMatch{}
	assert.False(router.Match(invalidMsg, nil, match))
}
