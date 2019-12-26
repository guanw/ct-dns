package http

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DecodeBody(t *testing.T) {
	tests := []struct {
		body        io.Reader
		expectedErr bool
		description string
		expected    postBody
	}{
		{
			body:        bytes.NewReader([]byte(`{random:"random"}`)),
			description: "empty body should throw error",
			expectedErr: true,
		},
		{
			body: bytes.NewReader([]byte(`{"serviceName":"dummy-service", "host":"192.0.0.1", "operation":"add"}`)),
			expected: postBody{
				ServiceName: "dummy-service",
				Host:        "192.0.0.1",
				Operation:   "add",
			},
			expectedErr: false,
			description: "body should be parsed correctly",
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			b, err := decodeBody(test.body)
			if test.expectedErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.expected, b)
			}
		})
	}
}
