package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	c := NewConfig()

	assert.NotEmpty(t, c)
	assert.True(t, c.IsSet("APP_ADDR"))
	assert.True(t, c.IsSet("SERVER_READ_TIMEOUT"))
	assert.True(t, c.IsSet("SERVER_READ_HEADER_TIMEOUT"))
	assert.True(t, c.IsSet("SERVER_WRITE_TIMEOUT"))
}
