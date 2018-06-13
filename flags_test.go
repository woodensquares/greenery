package greenery

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMisc(t *testing.T) {
	// Default out of valid range
	require.PanicsWithValue(t, "Invalid value -1 for variable , should be between 0 and 10", func() { NewDefaultIntValue("", -1, 0, 10) })
	require.PanicsWithValue(t, "While creating variable var, error nope", func() {
		NewDefaultCustomStringValue("var", "bad", func(n, s string, data interface{}) (string, error) {
			if s == "bad" {
				return "", fmt.Errorf("nope")
			}
			return s, nil
		}, nil)
	})
	require.PanicsWithValue(t, "While creating variable var, error nope", func() {
		NewCustomStringValue("var", func(n, s string, data interface{}) (string, error) {
			if s == "" {
				return "", fmt.Errorf("nope")
			}
			return s, nil
		}, nil)
	})

	// Default values
	fip := NewIPValue("hi")
	require.Equal(t, fip.Value, "0.0.0.0")
	fp := NewPortValue()
	require.EqualValues(t, fp, 0)

	// Unmarshaling
	fi := NewIntValue("hi", 2, 100)
	require.EqualValues(t, fi.Value, 2)
	require.NoError(t, fi.UnmarshalText([]byte("50")))
	require.EqualValues(t, fi.Value, 50)

	fc := NewCustomStringValue("", func(n, s string, data interface{}) (string, error) {
		return s, nil
	}, nil)
	require.NoError(t, fc.UnmarshalText([]byte("hi")))
	require.EqualValues(t, fc.Value, "hi")

	require.NoError(t, fp.UnmarshalText([]byte("")))
	require.EqualValues(t, fp, 0)
	require.Error(t, fp.UnmarshalText([]byte("hihihi")))

	require.NoError(t, fp.UnmarshalText([]byte("32")))
	require.Equal(t, fp.GetTyped(), uint16(32))

}
