package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)
type User struct{
	ID int64 `json:"id"`
	UserName string `json:"username"`
	Email string `json:"email"`
	Password password `json:"-"`
	CreatedAt string `json:"created_at"`
}

type password struct{
	text *string
	hash []byte
}

func (p *password) Set(text string) error{
	hash,err:=bcrypt.GenerateFromPassword([]byte(text),bcrypt.DefaultCost)
	if err!=nil{
		return err
	}
	p.text=&text
	p.hash=hash
	return nil
}

type UserStore struct{
	db *sql.DB
}
func (s *UserStore) Create(ctx context.Context,tx *sql.Tx,user *User) error{
	query:=`
	INSERT INTO users(username,email,password)
	VALUES($1,$2,$3) RETURNING id,created_at
	`
	ctx,cancel:=context.WithTimeout(ctx,QueryTimeoutDuration)
	defer cancel()
	err:=s.db.QueryRowContext(ctx, query,user.UserName,user.Email,user.Password).Scan(&user.ID,&user.CreatedAt)
	if err!=nil{
		switch {
		case err.Error()==`pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrorDuplicateEmail
		case err.Error()==`pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrorDuplicateUsername
		default:
			return err
		}
	}
	return nil
}

func (s *UserStore) GetById(ctx context.Context,id int64) (*User,error){
	query:=`
	SELECT id,username,email,created_at
	FROM users WHERE id=$1
	`
	var user User
	ctx,cancel:=context.WithTimeout(ctx,QueryTimeoutDuration)
	defer cancel()
	err:=s.db.QueryRowContext(ctx,query,id).Scan(&user.ID,&user.UserName,&user.Email,&user.CreatedAt)
	if err!=nil{
		switch {
		case	errors.Is(err,sql.ErrNoRows):
			return nil,ErrorNotFound
		default:
			return nil,err
		}
	}
	log.Printf("User: %v",user)
	return &user,err
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User,token string,invitationExp time.Duration) error{
	 return withTx(s.db,ctx,func(tx *sql.Tx) error{
		//create a user
		if err:=s.Create(ctx,tx,user);err!=nil{
			return err
		}
		//create invitation
		if err:=s.createUserInvitation(ctx,tx,token,invitationExp,user.ID);err!=nil{
			return err
		}
		return nil
	 })
}

func (s *UserStore) createUserInvitation(ctx context.Context,tx *sql.Tx,token string,invitationExp time.Duration,userId int64) error{
	query:=`INSERT INTO user_invitations (token,user_id,expiry) VALUES ($1,$2,$3)`
	ctx,cancel:=context.WithTimeout(ctx,QueryTimeoutDuration)
	defer cancel()
	_,err:=tx.ExecContext(ctx,query,token,userId,time.Now().Add(invitationExp))
	if err!=nil{
		return err
	} 
	return nil
}