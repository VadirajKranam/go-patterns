package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (app *application) getPostHandler(w http.ResponseWriter,r *http.Request){
	idParam:=chi.URLParam(r,"postId")
	id,err:=strconv.ParseInt(idParam,10,64)
	log.Printf("id: %v",id)
	if err!=nil{
		writeJSONError(w,r,http.StatusInternalServerError,err.Error())
	}
	ctx:=r.Context()
	post,err:=app.store.Posts.GetById(ctx,id)
	if err!=nil{
		switch{
		case errors.Is(err,store.ErrorNotFound):
			writeJSONError(w,r,http.StatusNotFound,err.Error())
		default:
			writeJSONError(w,r,http.StatusInternalServerError,err.Error())
		}
		
		return
	}
	if err:=writeJson(w,http.StatusOK,post);err!=nil{
		writeJSONError(w,r,http.StatusInternalServerError,err.Error())
		return
	}
}