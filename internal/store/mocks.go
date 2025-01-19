package store

import (
	"context"
	"database/sql"
	"time"
)

type MockUserStore struct{

}

func NewMockStore() Storage{
	return Storage{
		Users: &MockUserStore{},
	}
}

func (m *MockUserStore) Create(ctx context.Context,tx *sql.Tx,user *User) error{
	return nil
}

func (m *MockUserStore) GetById(ctx context.Context,userId int64) (*User,error){
	return &User{},nil
}

func (m *MockUserStore) CreateAndInvite(ctx context.Context, user *User,token string,invitationExp time.Duration) error{
	return nil
}

func (m *MockUserStore) createUserInvitation(ctx context.Context,tx *sql.Tx,token string,invitationExp time.Duration,userId int64) error{
	return nil
}
func (m *MockUserStore) Activate(context.Context,string) error{
	return nil
}

func (m *MockUserStore) Delete(ctx context.Context,userId int64) error{
	return nil
}

func (m *MockUserStore) GetByEmail(ctx context.Context,email string) (*User,error){
	return nil,nil
}