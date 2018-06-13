package greenery

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var bogus = struct {
	String string
	Int    int
}{String: "hi", Int: 12}

func TestGetCfgWrongStruct(t *testing.T) {
	base := NewBaseConfig("testing", nil)

	s, err := getCfg(&bogus)
	require.Nil(t, s, "Returned struct is not nil")
	require.Error(t, err, "Not a struct embedding BaseConfig", "Did not fail")

	s, err = getCfg(bogus)
	require.Nil(t, s, "Returned struct is not nil")
	require.Error(t, err, "Not a struct embedding BaseConfig", "Did not fail")

	s, err = getCfg(123)
	require.Nil(t, s, "Returned struct is not nil")
	require.Error(t, err, "Not a struct", "Did not fail")

	s, err = getCfg(base)
	require.Equal(t, base, s)
	require.NoError(t, err)

}
