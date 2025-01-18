package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vadiraj/gopher/internal/store"
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
			user,err:=app.getUser(ctx,userId)
			if err!=nil{
				app.unAuthorizedErrorResponse(w,r,err)
				return
			}
			ctx=context.WithValue(ctx,userCtx,user)
			next.ServeHTTP(w,r.WithContext(ctx))
		})
}


func (app *application) CheckPostOwnership(role string,next http.HandlerFunc) ( http.HandlerFunc){
	return http.HandlerFunc(func (w http.ResponseWriter,r *http.Request){
		user:=getUserFromCtx(r)
		post:=getPostFromCtx(r) 
		//if it is the user post
		if post.UserID==user.ID{
			next.ServeHTTP(w,r)
		}
		//role precedence check
		allowed,err:=app.checkRolePrecedence(r.Context(),user,role)
		if err!=nil{
			app.internalServerError(w,r,err)
			return
		}
		if !allowed{
			app.forbiddenResponse(w,r)
			return
		}	
		next.ServeHTTP(w,r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context,user *store.User,roleName string) (bool,error){
	role,err:=app.store.Roles.GetByName(ctx,roleName)
	if err!=nil{
		return false,err
	}
	return user.Role.Level>=role.Level,nil
}

func (app *application) getUser(ctx context.Context,userId int64) (*store.User,error){

	user,err:=app.cacheStorage.Users.Get(ctx,userId)
	if err!=nil{
		return nil,err
	}
	if user==nil{
	user,err=app.store.Users.GetById(ctx,userId)
		if err!=nil{
			return nil,err
		}
		if err:=app.cacheStorage.Users.Set(ctx,user);err!=nil{
			return nil,err
		}
	}
	return user,err
}