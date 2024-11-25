package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound     = errors.New("resource  not found")
	ErrConflict     = errors.New("resource already exists")
	TimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		Delete(ctx context.Context, id int64) error
		Update(context.Context, *Post) error
		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetadata, error)
	}
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetByID(context.Context, int64) (*User, error)
		CreateAndInvite(context.Context, *User, string, time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int64) error
	}
	Comments interface {
		GetByPostID(context.Context, int64) ([]Comment, error)
		Create(context.Context, *Comment) error
	}

	Followers interface {
		Follow(ctx context.Context, followerID, userID int64) error
		UnFollow(ctx context.Context, followerID, userID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostStore{
			db: db,
		},
		Users:     &UserStore{db: db},
		Comments:  &CommentStore{db: db},
		Followers: &FollowerStore{db: db},
	}
}

func withTx(db *sql.DB, ctx context.Context, f func(*sql.Tx) error) error {
	tx, _ := db.BeginTx(ctx, nil)
	if err := f(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
