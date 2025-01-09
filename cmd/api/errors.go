package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter,r *http.Request, err error){
	app.logger.Errorw("internal server error: ",r.Method, "path :",r.URL.Path,"error:",err)
	writeJSONError(w,http.StatusInternalServerError,"the server encouyntered a problem")
}

func (app *application) badRequestError(w http.ResponseWriter,r *http.Request,err error){
	app.logger.Warnf("bad request error: ",r.Method, "path :",r.URL.Path,"error:",err)
	writeJSONError(w,http.StatusBadRequest,err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter,r *http.Request,err error){
	app.logger.Warnf("not found error: ",r.Method, "path :",r.URL.Path,"error:",err)
	writeJSONError(w,http.StatusNotFound,"not found")
}

func (app *application) conflictResponse(w http.ResponseWriter,r *http.Request,err error){
	app.logger.Errorf("conflict error: ",r.Method, "path :",r.URL.Path,"error:",err)
	writeJSONError(w,http.StatusConflict,err.Error())
}