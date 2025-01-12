package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/vadiraj/gopher/internal/mailer"
	"github.com/vadiraj/gopher/internal/store"
)


type RegisterUserPayload struct{
	Username string `json:"username" validate:"required,max=100"`
	Email string `json:"email" validate:"required,email,max=225"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type UserWithToken struct{
	*store.User
	Token string `json:"token"`
}

func (app *application) registerUserHandler(w http.ResponseWriter,r *http.Request){
	var payLoad RegisterUserPayload
	if err:=readJson(w,r,&payLoad);err!=nil{
		app.badRequestError(w,r,err)
		return
	}
	if err:=Validate.Struct(payLoad);err!=nil{
		app.badRequestError(w,r,err)
		return
	}
	user:=&store.User{
		UserName: payLoad.Username,
		Email: payLoad.Email,
	}
	//hash the user password
	if err:=user.Password.Set(payLoad.Password);err!=nil{
		app.internalServerError(w,r,err)
		return
	}
	//store the user
	ctx:=r.Context()
	plainToken:=uuid.New().String()
	hash:=sha256.Sum256([]byte(plainToken))
	hashToken:=hex.EncodeToString(hash[:])
	err:=app.store.Users.CreateAndInvite(ctx,user,hashToken,app.config.mail.exp )
	if err!=nil{
		switch err{
		case store.ErrorDuplicateEmail:
			app.badRequestError(w,r,err)
		case store.ErrorDuplicateUsername:
			app.badRequestError(w,r,err)
		default:
			app.internalServerError(w,r,err)
		}
		return
	}

	userWithToken:=UserWithToken{
		User: user,
		Token: plainToken,
	}
	isProdEnv:=app.config.env=="production"
	activationUrl:=fmt.Sprintf("%s/confirm/%s",app.config.frontendUrl,plainToken)

	vars:=struct{
		Username string
		ActivationURL string
	}{
		Username: user.UserName,
		ActivationURL: activationUrl,
	}
	_,err=app.mailer.Send(mailer.UserWelcomeTemplate,user.UserName,user.Email,vars,!isProdEnv)
	if err!=nil{
		app.logger.Errorw("error sending the welcome email","error",err)
		//rollback user creation if email fails (SAGA pattern)
		if err:=app.store.Users.Delete(ctx,user.ID);err!=nil{
			app.logger.Errorw("error deleting user","err",err)
		}
		return
	}
	if err:=app.jsonResponse(w,http.StatusCreated,userWithToken);err!=nil{
		app.internalServerError(w,r,err)
		return 
	}
}