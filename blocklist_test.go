package layer

import (
	"testing"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/require"
)

func TestAddUserToBlockList(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()

	ok, err := l.AddUserToBlockList(user1, user2)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestGetUserBlockList(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()

	ok, err := l.AddUserToBlockList(user1, user2)
	require.NoError(t, err)
	require.True(t, ok)

	b, err := l.GetUserBlockList(user1)
	require.NoError(t, err)
	require.Len(t, b, 1)
	require.Equal(t, b[0].UserID, user2)
}

func TestUnblockUser(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()

	ok, err := l.AddUserToBlockList(user1, user2)
	require.NoError(t, err)
	require.True(t, ok)

	b, err := l.GetUserBlockList(user1)
	require.NoError(t, err)
	require.Len(t, b, 1)
	require.Equal(t, b[0].UserID, user2)

	ok, err = l.UnblockUser(user1, user2)
	require.NoError(t, err)
	require.True(t, ok)

	b, err = l.GetUserBlockList(user1)
	require.NoError(t, err)
	require.Len(t, b, 0)
}

func TestBulkModifyBlockList(t *testing.T) {

	user1 := uuid.New()
	user2 := uuid.New()
	user3 := uuid.New()

	ok, err := l.AddUserToBlockList(user1, user3)
	require.NoError(t, err)
	require.True(t, ok)

	b, err := l.GetUserBlockList(user1)
	require.NoError(t, err)
	require.Len(t, b, 1)
	require.Equal(t, b[0].UserID, user3)

	blockList := []string{user2}
	unBlockList := []string{user3}
	ok, err = l.BulkModifyBlockList(user1, blockList, unBlockList)
	require.NoError(t, err)
	require.True(t, ok)

	b, err = l.GetUserBlockList(user1)
	require.NoError(t, err)
	require.Len(t, b, 1)
	require.Equal(t, b[0].UserID, user2)

}
