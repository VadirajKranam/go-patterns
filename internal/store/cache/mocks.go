package cache

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/vadiraj/gopher/internal/store"
)

func NewMockStore() Storage{
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{
	mock.Mock
}

func (m *MockUserStore) Get(ctx context.Context,userId int64) (*store.User,error){
	args:=m.Called(userId)
	return nil,args.Error(1)
}

func (m *MockUserStore) Set(ctx context.Context,user *store.User) error{
	args:=m.Called(user)
	return args.Error(0)
}

func (m *MockUserStore) Delete(ctx context.Context,userId int64){
	m.Called(userId)
}