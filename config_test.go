package wsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckSet(t *testing.T) {
	for _, tc := range []struct {
		defaultName         string
		enable              []string
		disable             []string
		expectedChecks      CheckSet
		expectedErrContains string
	}{
		{
			defaultName:    "",
			expectedChecks: DefaultChecks(),
		},
		{
			defaultName:    "all",
			expectedChecks: AllChecks(),
		},
		{
			defaultName:    "none",
			expectedChecks: NoChecks(),
		},
		{
			defaultName:    "none",
			disable:        []string{"for", "range"},
			expectedChecks: NoChecks(),
		},
		{
			defaultName: "none",
			enable:      []string{"for", "range"},
			disable:     []string{"range"},
			expectedChecks: map[CheckType]struct{}{
				CheckFor: {},
			},
		},
		{
			defaultName: "none",
			enable:      []string{"for", "range"},
			expectedChecks: map[CheckType]struct{}{
				CheckFor:   {},
				CheckRange: {},
			},
		},
		{
			defaultName:         "all",
			disable:             []string{"invalid-disable"},
			expectedErrContains: "invalid check 'invalid-disable'",
		},
		{
			defaultName:         "all",
			enable:              []string{"invalid-enable"},
			expectedErrContains: "invalid check 'invalid-enable'",
		},
		{
			defaultName:         "invalid",
			expectedErrContains: "invalid preset",
		},
	} {
		t.Run(tc.defaultName, func(t *testing.T) {
			checks, err := NewCheckSet(tc.defaultName, tc.enable, tc.disable)
			if tc.expectedErrContains != "" {
				assert.Contains(t, err.Error(), tc.expectedErrContains)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedChecks, checks)
		})
	}
}

func TestToAndFromString(t *testing.T) {
	maxCheckNumber := 23

	for n := range maxCheckNumber {
		check := CheckType(n)
		ct, err := CheckFromString(check.String())

		if n == 0 {
			assert.Equal(t, "invalid", check.String())
			require.Error(t, err)

			continue
		}

		require.NoError(t, err)
		assert.Equal(t, check, ct)
	}
}
