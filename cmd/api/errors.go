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

func (app *application) unAuthorizedErrorResponse(w http.ResponseWriter,r *http.Request,err error){
	app.logger.Errorf("unauthorized error: ",r.Method, "path :",r.URL.Path,"error:",err)
	writeJSONError(w,http.StatusUnauthorized,err.Error())
}

func (app *application) unAuthorizedBasicErrorResponse(w http.ResponseWriter,r *http.Request,err error){
	app.logger.Warnf("unauthorized basic error: ",r.Method, "path :",r.URL.Path,"error:",err)
	w.Header().Set("WWW-Authenticate",`Basic realm="restricted",charset="UTF-8"`)
	writeJSONError(w,http.StatusUnauthorized,err.Error())
}
