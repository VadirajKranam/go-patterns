package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vadiraj/gopher/internal/auth"
	"github.com/vadiraj/gopher/internal/store"
	"github.com/vadiraj/gopher/internal/store/cache"
	"go.uber.org/zap"
)


func newTestApplication(t *testing.T,config config) *application{
	t.Helper()
	logger:=zap.NewNop().Sugar()
	//logger := zap.Must(zap.NewProduction()).Sugar()
	mockStore:=store.NewMockStore()
	mockCacheStore:=cache.NewMockStore()
	testAuth:=&auth.TestAuthenticator{}
	return &application{
		logger: logger,
		store:mockStore,
		authenticator: testAuth,
		cacheStorage: mockCacheStore,
		config: config,
	}
}
func execRequest(req *http.Request,mux http.Handler) *httptest.ResponseRecorder{
	rr:=httptest.NewRecorder()
	mux.ServeHTTP(rr,req)
	return rr
}

func checkResponseCode(t *testing.T, expected,actual int){
	if expected!=actual{
		t.Errorf("expected the response code to be %d and we got %d",expected,actual)
	}
}