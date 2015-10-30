package mux

import (
	"testing"

	"github.com/stretchr/testify/assert"

	coap "github.com/dustin/go-coap"
)

func TestPathMatcherSingle(t *testing.T) {
	assert := assert.New(t)
	pathRoute := &Route{}
	pathRoute.Path("/stats/{id}/clients")

	testMsg := &coap.Message{}
	testMsg.SetPathString("/stats/1234/clients")
	match := &RouteMatch{}
	assert.True(pathRoute.Match(testMsg, nil, match))
	assert.Equal(1, len(match.Vars))
}

func TestPathMatcherMultiple(t *testing.T) {
	assert := assert.New(t)
	pathRoute := &Route{}
	pathRoute.Path("/stats/{id}/{subid}/clients")

	testMsg := &coap.Message{}
	testMsg.SetPathString("/stats/1234/abcd/clients")
	match := &RouteMatch{}
	assert.True(pathRoute.Match(testMsg, nil, match))
	assert.Equal(2, len(match.Vars))
}
