package phone_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/phone"
)

func TestNormalize(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"09120000000", "+989120000000"},
		{"+989120000000", "+989120000000"},
	}
	for _, tc := range cases {
		got, err := phone.Normalize(tc.in)
		require.NoError(t, err)
		assert.Equal(t, tc.out, got)
	}
}
