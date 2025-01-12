package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/vadiraj/gopher/internal/store"
)

type userKey string

const userCtx userKey = "user"

// GetUser godoc
// @Summary      Fetches a user profile
// @Description  Fetches a user profile by id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  store.User
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /user/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

type FollowUser struct {
	UserId int64 `json:"user_id"`
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)
	ctx := r.Context()

	//Revert back to auth userID from ctx
	var payLoad FollowUser
	if err := readJson(w, r, &payLoad); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	log.Print("followerid: ", followerUser.ID, "userid :", payLoad.UserId)
	if err := app.store.Followers.Follow(ctx, followerUser.ID, payLoad.UserId); err != nil {
		switch {
		case errors.Is(err, store.ErrorConflict):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)
	ctx := r.Context()

	//Revert back to auth userID from ctx
	var payLoad FollowUser
	if err := readJson(w, r, payLoad); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := app.store.Followers.Unfollow(ctx, followerUser.ID, payLoad.UserId); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func(app *application) activateUserHandler(w http.ResponseWriter,r *http.Request){
	token:=chi.URLParam(r,"token")
	err:=app.store.Users.Activate(r.Context(),token)
	if err!=nil{
		switch err{
		case store.ErrorNotFound:
			app.notFoundError(w,r,err)
		default:
			app.internalServerError(w,r,err)
		}
		return
	}
	if err:=app.jsonResponse(w,http.StatusNoContent,"");err!=nil{
		app.internalServerError(w,r,err)
		return
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userId")
		userId, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
		}
		ctx := r.Context()
		user, err := app.store.Users.GetById(ctx, userId)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrorNotFound):
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
		}
		log.Printf("userId: %v", userId)
		ctx = context.WithValue(r.Context(), userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
