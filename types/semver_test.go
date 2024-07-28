package types_test

import (
	"testing"

	"github.com/ffenix113/zigbee_home/types"
	"github.com/stretchr/testify/require"
)

func TestSemverCompare(t *testing.T) {
	tests := []struct {
		name       string
		ver1, ver2 string
		result     int
	}{
		{
			name:   "equal",
			ver1:   "v1.1.1",
			ver2:   "v1.1.1",
			result: 0,
		},
		{
			name:   "less",
			ver1:   "v1.6.200",
			ver2:   "v1.7.0",
			result: -1,
		},
		{
			name:   "greater",
			ver1:   "v1.0.0",
			ver2:   "v0.3",
			result: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ver1, err := types.ParseSemver(test.ver1)
			require.NoError(t, err)

			ver2, err := types.ParseSemver(test.ver2)
			require.NoError(t, err)

			require.Equal(t, test.result, ver1.Compare(ver2))
		})
	}
}
