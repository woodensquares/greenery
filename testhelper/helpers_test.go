package testhelper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompareIgnoreTmp(t *testing.T) {
	file1 := []byte(`some line
some /tmp/file1 line
/tmp/ignorethat line
line /tmp/file1
but /tmp/not!a!tmp file`)

	file2 := []byte(`some line
some /tmp/file2 line
/tmp/ignorethat line
line /tmp/file2
but /tmp/not!a!tmp file`)

	file3 := []byte(`some line
some /tmp/file3 line
/tmp/ignorethat line
line /tmp/file3
but /tmp/not!a!different!tmp file`)

	file4 := []byte(`some line
some /tmp/file4 line
/tmp/ignorethat line
line /tmp/file4
/tmp/ignorethat line
line /tmp/file4
/tmp/ignorethat line
line /tmp/file4
but /tmp/not!a!tmp file`)

	require.True(t, CompareIgnoreTmp(t, file1, file2))
	require.False(t, CompareIgnoreTmp(t, file1, file3))
	require.False(t, CompareIgnoreTmp(t, file1, file4))

}
