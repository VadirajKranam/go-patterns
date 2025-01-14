package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (app *application) BasicAuthMiddleware(next http.Handler) (http.Handler){
	
		return	http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
			//read the auth header
			authHeader:=r.Header.Get("Authorization")
			if authHeader==""{
				app.unAuthorizedBasicErrorResponse(w,r,fmt.Errorf("authorization header is missing"))
				return
			}
			// parse it-> get the base64
			parts:=strings.Split(authHeader," ")
			if len(parts)!=2 || parts[0]!="Basic"{
				app.unAuthorizedBasicErrorResponse(w,r,fmt.Errorf("authorization header is ,alformed"))
				return
			}
			// decode it
			decoded,err:=base64.StdEncoding.DecodeString(parts[1])
			if err!=nil{
				app.unAuthorizedBasicErrorResponse(w,r,fmt.Errorf("authorization header is ,alformed"))
				return
			}
			username:=app.config.auth.basic.user
			pass:=app.config.auth.basic.pass
			// check the credentials
			creds:=strings.SplitN(string(decoded),":",2)
			if len(creds)!=2 || creds[0]!=username || creds[1]!=pass{
				app.unAuthorizedBasicErrorResponse(w,r,fmt.Errorf("invalid credentials"))
				return
			}
			next.ServeHTTP(w,r)
		})
}

func (app *application) AuthTokenMiddleware(next http.Handler) ( http.Handler){
		return http.HandlerFunc(func (w http.ResponseWriter,r *http.Request){
			authHeader:=r.Header.Get("Authorization")
			if authHeader==""{
				app.unAuthorizedErrorResponse(w,r,fmt.Errorf("authorization header is missing"))
				return
			}
			parts:=strings.Split(authHeader," ")
			if len(parts)!=2 || parts[0]!="Bearer"{
				app.unAuthorizedErrorResponse(w,r,fmt.Errorf("authorization header is ,alformed"))
				return
			}
			token:=parts[1]
			jwtToken,err:=app.authenticator.ValidateToken(token)
			if err!=nil{
				app.unAuthorizedErrorResponse(w,r,err)
				return
			}
			claims:=jwtToken.Claims.(jwt.MapClaims)
			userId,err:=strconv.ParseInt(fmt.Sprintf("%.f",claims["sub"]),10,64)
			if err!=nil{
				app.unAuthorizedErrorResponse(w,r,err)
				return
			}
			ctx:=r.Context()
			user,err:=app.store.Users.GetById(ctx,userId)
			if err!=nil{
				app.unAuthorizedErrorResponse(w,r,err)
				return
			}
			ctx=context.WithValue(ctx,userCtx,user)
			next.ServeHTTP(w,r.WithContext(ctx))
		})
}