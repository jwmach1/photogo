package photos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_buildPath(t *testing.T) {
	ct, err := time.Parse(time.RFC3339, "1994-05-20T15:04:05Z")
	require.NoError(t, err)

	actual := buildPath(("base"), ct)

	expected := "base/1994/05"
	assert.Equal(t, expected, actual)
}

func Test_(t *testing.T) {

}
