package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrorNotFound=errors.New("resource not found")
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
	}
}
func NewStorage(db *sql.DB) Storage{
	return Storage{
		Posts: &PostStore{db:db},
		Users: &UserStore{db:db},
		Comments: &CommentStore{db: db},
	}
}