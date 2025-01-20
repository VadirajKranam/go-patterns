package main

import (
	"net/http"

	"github.com/vadiraj/gopher/internal/store"
)

type AddCommentPayload struct {
	Content string `json:"content" validate:"required,max=100"`
}

func (app *application) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payLoad AddCommentPayload
	if err := readJson(w, r, &payLoad); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payLoad); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	post := getPostFromCtx(r)
	ctx := r.Context()
	comment := &store.Comment{
		Content: payLoad.Content,
		PostID:  post.ID,
		//todo remove after auth
		UserID: 1,
	}
	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
