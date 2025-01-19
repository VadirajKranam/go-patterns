package cache

import (
	"context"
	"log"

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
	log.Printf("Mock Get called with: ctx=%v, userId=%v", ctx, userId)
	args:=m.Called(ctx,userId)
	return nil,args.Error(1)
}

func (m *MockUserStore) Set(ctx context.Context,user *store.User) error{
	args:=m.Called(ctx,user)
	return args.Error(0)
}

func (m *MockUserStore) Delete(ctx context.Context,userId int64){
	m.Called(ctx,userId)
}