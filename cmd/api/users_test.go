package main

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/vadiraj/gopher/internal/store/cache"
)

func TestGetUser(t *testing.T) {
	withRedis := config{
		redisCfg: redisConfig{
			enabled: true,
		},
	}
	app := newTestApplication(t, withRedis)
	log.Print("enabled: ", app.config.redisCfg.enabled)
	mux := app.mount()
	testToken, err := app.authenticator.GenerateToken(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		//check for 401 code
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := execRequest(req, mux)
		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})
	t.Run("should allow authenticated requests", func(t *testing.T) {
		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserStore)
		mockCacheStore.On("Get", mock.Anything, int64(1)).Return(nil, nil)
		mockCacheStore.On("Set", mock.Anything).Return(nil)
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+testToken)
		rr := execRequest(req, mux)
		checkResponseCode(t, http.StatusOK, rr.Code)
		mockCacheStore.AssertNumberOfCalls(t, "Get", 1)
		mockCacheStore.AssertCalled(t, "Set", mock.Anything)
		mockCacheStore.Calls = nil
	})
	t.Run("should hit the cache first and if not exists it sets the user on the cache", func(t *testing.T) {
		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserStore)
		mockCacheStore.On("Get", mock.Anything, int64(42)).Return(nil, nil)
		mockCacheStore.On("Get", mock.Anything, int64(1)).Return(nil, nil)
		mockCacheStore.On("Set", mock.Anything, mock.Anything).Return(nil)
		//	mockCacheStore.On("Delete",mock.Anything,mock.Anything)
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+testToken)
		rr := execRequest(req, mux)
		checkResponseCode(t, http.StatusOK, rr.Code)
		mockCacheStore.AssertNumberOfCalls(t, "Get", 2)
		mockCacheStore.Calls = nil
	})
}
