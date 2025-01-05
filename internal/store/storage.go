package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrorNotFound=errors.New("resource not found")
	QueryTimeoutDuration=time.Second*5
)

type Storage struct{
	Posts interface{
		Create(context.Context,*Post) error
		GetById(context.Context,int64) (*Post,error)
		Delete(context.Context,int64) error
		Update(context.Context,*Post) error
	}
	Users interface{
		Create(context.Context,*User) error
	}
	Comments interface{
		GetByPostID(context.Context,int64) ([]Comment,error)
		Create(context.Context,*Comment) error
	}
}
func NewStorage(db *sql.DB) Storage{
	return Storage{
		Posts: &PostStore{db:db},
		Users: &UserStore{db:db},
		Comments: &CommentStore{db: db},
	}
}