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
	ErrorDuplicateEmail=errors.New("email already exists")
	ErrorDuplicateUsername=errors.New("username already exists")
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
		Create(context.Context,*sql.Tx,*User) error
		GetById(context.Context,int64) (*User,error)
		CreateAndInvite(ctx context.Context, user *User,token string,invitationExp time.Duration) error
		createUserInvitation(ctx context.Context,tx *sql.Tx,token string,invitationExp time.Duration,userId int64) error
		Activate(context.Context,string) error
		Delete(ctx context.Context,userId int64) error
		GetByEmail(ctx context.Context,email string) (*User,error)
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

func withTx(db *sql.DB,ctx context.Context,fn func(*sql.Tx) error) error{
	tx,err:=db.BeginTx(ctx,nil)
	if err!=nil{
		return err
	}
	if err:=fn(tx);err!=nil{
		_=tx.Rollback()
		return err
	}
	return tx.Commit()
}