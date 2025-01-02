package main

import (
	"net/http"

	"github.com/vadiraj/gopher/internal/store"
)

type CreatePostPayload struct{
	Title string `json:"title"`
	Content string `json:"content"`
	Tags []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter,r *http.Request){
	var payLoad CreatePostPayload
	if err:=readJson(w,r,&payLoad);err!=nil{
		writeJSONError(w,r,http.StatusBadRequest,err.Error())
		return
	}
	post:=&store.Post{
		Title: payLoad.Title,
		Content: payLoad.Content,
		Tags: payLoad.Tags,
		//todo change after auth
		UserID: 1,
	}
	ctx:=r.Context()
	if err:=app.store.Posts.Create(ctx,post);err!=nil{
		writeJSONError(w,r,http.StatusInternalServerError,err.Error())
		return
	}
	if err:=writeJson(w,http.StatusCreated,post);err!=nil{
		writeJSONError(w,r,http.StatusInternalServerError,err.Error())
		return
	}
}