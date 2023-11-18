package db

import (
	"awesomeProject/utils"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomPost(t *testing.T) Post {
	firstAccount := createRandomAccount(t)
	secondAccount := createRandomAccount(t)

	args := CreatePostParams{
		FromAccountID: firstAccount.ID,
		ToAccountID:   secondAccount.ID,
		Amount:        utils.RandomMoney(),
	}

	post, err := testQueries.CreatePost(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, post)

	require.Equal(t, post.Amount, args.Amount)
	require.Equal(t, post.ToAccountID, args.ToAccountID)
	require.Equal(t, post.FromAccountID, args.FromAccountID)

	require.NotZero(t, post.ID)
	require.NotZero(t, post.CreatedAt)

	return post
}

func TestCreatePost(t *testing.T) {
	createRandomPost(t)
}

func TestGetPost(t *testing.T) {
	post := createRandomPost(t)

	postFromDB, err := testQueries.GetPost(context.Background(), post.ID)

	require.NoError(t, err)
	require.NotEmpty(t, postFromDB)

	require.Equal(t, post.ID, postFromDB.ID)
	require.Equal(t, post.Amount, postFromDB.Amount)
	require.Equal(t, post.ToAccountID, postFromDB.ToAccountID)
	require.Equal(t, post.FromAccountID, postFromDB.FromAccountID)
	require.WithinDuration(t, post.CreatedAt, postFromDB.CreatedAt, time.Second)
}

func TestGetPosts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomPost(t)
	}

	args := GetPostsParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.GetPosts(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, post := range entries {
		require.NotEmpty(t, post)
	}
}

func TestUpdatePost(t *testing.T) {
	post := createRandomPost(t)

	args := UpdatePostParams{
		Amount:        utils.RandomMoney(),
		ToAccountID:   post.FromAccountID,
		FromAccountID: post.ToAccountID,
		ID:            post.ID,
	}

	updatedPost, err := testQueries.UpdatePost(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, updatedPost)

	require.Equal(t, post.ID, updatedPost.ID)
	require.Equal(t, args.Amount, updatedPost.Amount)
	require.Equal(t, args.ToAccountID, updatedPost.ToAccountID)
	require.Equal(t, args.FromAccountID, updatedPost.FromAccountID)
	require.WithinDuration(t, post.CreatedAt, updatedPost.CreatedAt, time.Second)

}

func TestDeletePost(t *testing.T) {
	post := createRandomPost(t)

	err := testQueries.DeletePost(context.Background(), post.ID)

	require.NoError(t, err)

	post, err = testQueries.GetPost(context.Background(), post.ID)

	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, post)
}
