package photos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_buildURL(t *testing.T) {
	t.Run("video", func(t *testing.T) {
		actual := buildURL("video/mpeg", "base")

		assert.Equal(t, "base=dv", actual)
	})
	t.Run("image", func(t *testing.T) {
		actual := buildURL("image/jpeg", "foo")

		assert.Equal(t, "foo=d", actual)
	})
}

func Test_buildPath(t *testing.T) {
	ct, err := time.Parse(time.RFC3339, "1994-05-20T15:04:05Z")
	require.NoError(t, err)

	actual := buildPath(("base"), ct)

	expected := "base/1994/05"
	assert.Equal(t, expected, actual)
}
