package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/vadiraj/gopher/docs" //generate swagger doc
	"github.com/vadiraj/gopher/internal/mailer"
	"github.com/vadiraj/gopher/internal/store"
	"go.uber.org/zap"
)

type application struct{
	config config
	store  store.Storage
	logger *zap.SugaredLogger
	mailer mailer.Client
}

type mailConfig struct{
	sendGrid sendGridConfig
	exp time.Duration
	fromEmail string
	mailTrap mailTrapConfig
}

type sendGridConfig struct{
	apiKey string
}

type mailTrapConfig struct{
	apiKey string
}

type config struct{
	addr string
	db dbConfig
	env string
	apiURL string
	mail mailConfig
	frontendUrl string	
}

type dbConfig  struct{
	addr string
	maxOpenConnns int
	maxIdleConns int
	maxIdleTime string
}

func (app *application) mount() http.Handler{
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60*time.Second))
	r.Route("/v1",func(r chi.Router){
		r.Get("/health", app.healthCheckHandler)
		docsUrl:=fmt.Sprintf("%s/swagger/doc.json",app.config.addr)
		r.Get("/swagger/*",httpSwagger.Handler(httpSwagger.URL(docsUrl)))
		r.Route("/posts",func(r chi.Router){
			r.Post("/",app.createPostHandler)
			r.Route("/{postId}",func(r chi.Router){
				r.Use(app.postsContextMiddleware)
				r.Get("/",app.getPostHandler)
				r.Patch("/",app.updatePostHandler)
				r.Delete("/",app.deletePostHandler)
				r.Post("/comment",app.addCommentHandler)
			})
		})
		r.Route("/users",func(r chi.Router){
			r.Put("/activate/{token}",app.activateUserHandler)
			r.Route("/{userId}",func(r chi.Router){
				r.Use(app.userContextMiddleware)
				r.Get("/",app.getUserHandler)
				r.Put("/follow",app.followUserHandler)
				r.Put("/unfollow",app.unfollowUserHandler)
			})
			r.Group(func (r chi.Router)  {
				r.Get("/feed",app.getUserFeedHandler)
			})
		})
		//Public routes
		r.Route("/authentication",func(r chi.Router){
			r.Post("/user",app.registerUserHandler)
		})
	})

	return r
}
func (app *application) run(mux http.Handler) error{
	docs.SwaggerInfo.Version=version
	docs.SwaggerInfo.Host=app.config.apiURL
	docs.SwaggerInfo.BasePath="/v1"

	srv:=&http.Server{
		Addr: app.config.addr,
		Handler: mux,
		WriteTimeout: time.Second*30,
		ReadTimeout: time.Second*10,
		IdleTimeout: time.Minute,
	}
	app.logger.Infow("server has started at",app.config.addr ,"env: ",app.config.env)
	return srv.ListenAndServe()
}