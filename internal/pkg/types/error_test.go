package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorVariables(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		expectedError string
	}{
		{"ErrNoObjectFound", ErrNoObjectFound, "no object found"},
		{"ErrNoBucketFound", ErrNoBucketFound, "no bucket found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualError(t, tt.err, tt.expectedError)
		})
	}
}
