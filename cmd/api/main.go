package main

import (
	"log"
	"time"

	"github.com/vadiraj/gopher/internal/db"
	"github.com/vadiraj/gopher/internal/env"
	"github.com/vadiraj/gopher/internal/mailer"
	"github.com/vadiraj/gopher/internal/store"
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
		},
	}
	//logger 
	logger:=zap.Must(zap.NewProduction()).Sugar()
	//database
	db,err:=db.New(cfg.db.addr, cfg.db.maxOpenConnns,cfg.db.maxIdleConns,cfg.db.maxIdleTime)
	logger.Info("Database connection pool established")
	if err!=nil{
		logger.Fatal(err)
	}
	defer db.Close()
	store:=store.NewStorage(db)
	//mailer:=mailer.NewSendGrid(cfg.mail.sendGrid.apiKey,cfg.mail.fromEmail)
	mailer,err:=mailer.NewMailTrap(cfg.mail.mailTrap.username,cfg.mail.mailTrap.password,cfg.mail.fromEmail)
	if err!=nil{
		log.Fatal(err)
	}
	app:=&application{
		config: cfg,
		store: store,
		logger:logger,
		mailer: mailer,
	}
	mux:=app.mount()
	logger.Fatal(app.run(mux))
}
 