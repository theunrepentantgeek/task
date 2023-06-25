package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateId(t *testing.T) {
	var b graphBuilder
	cases := []struct {
		id       string
		expected string
	}{
		{"name", "name"},
		{"first second", "first_second"},
		{"first:second", "first_second"},
		{"first:second:third", "first_second_third"},
		{":name", "name"},
		{"name:", "name"},
	}

	for _, c := range cases {
		c := c
		t.Run(c.id, func(t *testing.T) {
			t.Parallel()
			actual := b.createId(c.id)
			assert.Equal(t, c.expected, actual)
		})
	}
}
