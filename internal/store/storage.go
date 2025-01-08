package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrorNotFound=errors.New("resource not found")
	ErrorConflict=errors.New("resource already exists")
	QueryTimeoutDuration=time.Second*5
)

type Storage struct{
	Posts interface{
		Create(context.Context,*Post) error
		GetById(context.Context,int64) (*Post,error)
		Delete(context.Context,int64) error
		Update(context.Context,*Post) error
		GetUserFeed(context.Context,int64,PaginatedFeedQuery) ([]PostWithMetadata,error)
	}
	Users interface{
		Create(context.Context,*User) error
		GetById(context.Context,int64) (*User,error)

	}
	Comments interface{
		GetByPostID(context.Context,int64) ([]Comment,error)
		Create(context.Context,*Comment) error
	}
	Followers interface{
		Follow(context.Context, int64,int64) error
		Unfollow(context.Context, int64,int64) error
	}
}
func NewStorage(db *sql.DB) Storage{
	return Storage{
		Posts: &PostStore{db:db},
		Users: &UserStore{db:db},
		Comments: &CommentStore{db: db},
		Followers: &FollowerStore{db: db},
	}
}