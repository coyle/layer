package layer

import (
	"testing"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/require"
)

func TestSetUsersBadge(t *testing.T) {
	user1 := uuid.New()

	ok, err := l.SetUsersBadge(user1, 12)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestGetUsersBadge(t *testing.T) {
	user1 := uuid.New()

	ok, err := l.SetUsersBadge(user1, 12)
	require.NoError(t, err)
	require.True(t, ok)

	res, err := l.GetUsersBadge(user1)
	require.NoError(t, err)

	require.Equal(t, res.UnreadExternal, 12)
}
