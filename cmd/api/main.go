package main

import (
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vadiraj/gopher/internal/auth"
	"github.com/vadiraj/gopher/internal/db"
	"github.com/vadiraj/gopher/internal/env"
	"github.com/vadiraj/gopher/internal/mailer"
	"github.com/vadiraj/gopher/internal/store"
	"github.com/vadiraj/gopher/internal/store/cache"
	"go.uber.org/zap"
)

const version="0.0.1"

//	@title			Gopher social
//	@description	API for GopherSocial, a social network fo gophers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html


//	@BasePath	/v1
//@securityDefinitions.apiKey->ApiKeyAuth
//@in->->->->->->->header
//@name->->->->->->Authorization
//@description 
func main(){
	cfg:=config{
		addr: env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL","localhost:8080"),
		frontendUrl: env.GetString("FRONTEND_URL","localhost:4000"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "postgres://admin:admin@localhost/socialnetwork?sslmode=disable"),
			maxOpenConnns:  env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS",30),
			maxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisCfg: redisConfig{
			addr: env.GetString("REDIS_ADDR","localhost:6379"),
			pw: env.GetString("REDIS_PW",""),
			db: env.GetInt("REDIS_DB",0),
			enabled: env.GetBool("REDIS_ENABLED",false),
		},
		env: env.GetString("ENV","development"),
		mail: mailConfig{
			exp: time.Hour*24*3,//3 days
			fromEmail: env.GetString("FROM_EMAIL",""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY",""),
			},
			mailTrap:mailTrapConfig{
				username: env.GetString("MAILTRAP_USERNAME",""),
				password: env.GetString("MAILTRAP_PASSWORD",""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER","admin"),
				pass:env.GetString("AUTH_BASIC_PASSWORD","admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET","example"),
				exp: time.Hour*24*3,
				iss: "gophersocial",
			},
		},
	}
	//logger 
	logger:=zap.Must(zap.NewProduction()).Sugar()
	//database
	db,err:=db.New(cfg.db.addr, cfg.db.maxOpenConnns,cfg.db.maxIdleConns,cfg.db.maxIdleTime)
	if err!=nil{
		logger.Fatal(err)
	}
	logger.Info("Database connection pool established")
	var rdb *redis.Client
	if cfg.redisCfg.enabled{
		rdb=cache.NewRedisClient(cfg.redisCfg.addr,cfg.redisCfg.pw,cfg.redisCfg.db)
		logger.Info("redis connection established")
	}
	defer db.Close()
	store:=store.NewStorage(db)
	cacheStorage:=cache.NewRedisStorage(rdb)
	//mailer:=mailer.NewSendGrid(cfg.mail.sendGrid.apiKey,cfg.mail.fromEmail)
	mailer,err:=mailer.NewMailTrap(cfg.mail.mailTrap.username,cfg.mail.mailTrap.password,cfg.mail.fromEmail)
	if err!=nil{
		log.Fatal(err)
	}
	jwtAuthenticator:=auth.NewJWTAuthenticator(cfg.auth.token.secret,cfg.auth.token.iss,cfg.auth.token.iss)
	app:=&application{
		config: cfg,
		store: store,
		cacheStorage:*cacheStorage,
		logger:logger,
		mailer: mailer,
		authenticator: jwtAuthenticator,
	}
	mux:=app.mount()
	logger.Fatal(app.run(mux))
}
 