package log

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLevel(t *testing.T) {
	var l Level
	require.Equal(t, "error", l.String())

	SetLogLevel(LogWarning)
	require.Equal(t, LogWarning, _lg.level)
	SetLogLevel(LogErr)
}
