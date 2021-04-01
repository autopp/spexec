package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestGetName(t *testing.T) {
	cases := []struct {
		name     string
		test     *Test
		expected string
	}{
		{
			name:     "Name is not empty",
			test:     &Test{Name: "test of echo", Command: []string{"echo", "hello"}},
			expected: "test of echo",
		},
		{
			name:     "Name is empty",
			test:     &Test{Name: "", Command: []string{"echo", "hello"}},
			expected: "echo hello",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.test.GetName())
		})
	}
}
