package types_test

import (
	"testing"

	"github.com/ffenix113/zigbee_home/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestPin(t *testing.T) {
	t.Run("unmarshal yaml", func(t *testing.T) {
		type s struct {
			Pin types.Pin
		}

		tests := []struct {
			name        string
			skipTest    bool
			raw         string
			unmarshaled types.Pin
			err         bool
		}{
			{
				name: "valid short",
				raw:  `pin: 0.00`,
				unmarshaled: types.Pin{
					Port: types.NewOption(uint8(0)),
					Pin:  types.NewOption(uint8(0)),
				},
			},
			{
				// FIXME: Unmarshaler defined on Pin
				// will not be called for null yaml values.
				// This results in unmarshaling passing
				// and this test becoming worthless.
				// Assume to be non-critical issue,
				// and to be mitigated with configuration
				// validations.
				name: "invalid short, no pins",
				raw:  `pin: null`,
				// Special case, to not assume that the test is passing.
				skipTest: true,
			},
			{
				name: "invalid long, no pins",
				raw:  `pin: {}`,
				err:  true,
			},
			{
				name: "valid long",
				raw:  `pin: {port: 1, pin: 03, inverted: true}`,
				unmarshaled: types.Pin{
					Port:     types.NewOption(uint8(1)),
					Pin:      types.NewOption(uint8(3)),
					Inverted: true,
				},
			},
			{
				name: "invalid long",
				raw:  `pin: {port: 5, pin: 03, inverted: true}`,
				err:  true,
			},
			{
				name: "valid long, no pins",
				raw:  `pin: {id: pin1}`,
				unmarshaled: types.Pin{
					ID:   "pin1",
					Port: types.NewEmptyOption[uint8](),
					Pin:  types.NewEmptyOption[uint8](),
				},
			},
		}

		for _, test := range tests {
			test := test

			t.Run(test.name, func(t *testing.T) {
				if test.skipTest {
					t.SkipNow()
				}

				var s struct {
					Pin types.Pin
				}

				err := yaml.Unmarshal([]byte(test.raw), &s)
				if test.err {
					// TODO: We can check that error matches as well.
					require.Error(t, err)
					return
				}

				require.NoError(t, err)
				require.Equal(t, test.unmarshaled, s.Pin)
			})
		}
	})
}
