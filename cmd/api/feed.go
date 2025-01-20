package main

import (
	"net/http"

	"github.com/vadiraj/gopher/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	//pagination,filters
	ctx := r.Context()
	fq := store.PaginatedFeedQuery{
		Limit:  10,
		Offset: 0,
		Sort:   "asc",
	}
	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
	}
	if err := Validate.Struct(fq); err != nil {
		app.badRequestError(w, r, err)
	}

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(108), fq)
	if err != nil {
		app.internalServerError(w, r, err)
	}
	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}

}
