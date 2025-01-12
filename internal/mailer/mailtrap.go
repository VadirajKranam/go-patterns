package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"time"

	gomail "gopkg.in/mail.v2"
)


type mailTrapClient struct{
	fromEmail string
	username string
	password string
}

func NewMailTrap(username,password,fromEmail string) (*mailTrapClient,error){
	if username==""{
		return &mailTrapClient{},errors.New("username is required")
	}
	if password==""{
		return &mailTrapClient{},errors.New("password is required")
	}
	return &mailTrapClient{
		fromEmail: fromEmail,
		username: username,
		password: password,
	},nil
}

func (m *mailTrapClient) Send(templateFile,username,email string,data any,isSandbox bool) (int,error){
	entries, err := FS.ReadDir("template")
	if err != nil {
    		log.Printf("Error reading embedded template dir: %v", err)
	}
	for _, entry := range entries {
    	log.Printf("Found embedded file: %s", entry.Name())
	}	
	tmpl,err:=template.ParseFS(FS,"template/"+templateFile)
	if err!=nil{
		return -1,err
	}
	subject:=new(bytes.Buffer)
	err=tmpl.ExecuteTemplate(subject,"subject",data)
	if err!=nil{
		return -1,err
	}
	body:=new(bytes.Buffer)
	err=tmpl.ExecuteTemplate(body,"body",data)
	if err!=nil{
		return -1,err
	}
	var retryError error
	for i:=0;i<maxRetries;i++{
	message:=gomail.NewMessage()
	message.SetHeader("From",m.fromEmail)
	message.SetHeader("To",email)
	message.SetHeader("Subject",subject.String())
	message.AddAlternative("text/html",body.String())
	dialer:=gomail.NewDialer("smtp.mailosaur.net",587,m.username,m.password)
	if retryError=dialer.DialAndSend(message);retryError!=nil{
		log.Printf("Failed to send email to %v, attempt %d of %d",email,i+1,maxRetries)
			log.Printf("Error %v",retryError.Error())
			//exponential backoff
			time.Sleep(time.Second*time.Duration(i+1))
			continue
	}
	log.Printf("Email sent with status code: %v",200)
		return 200,nil
	}	
	return 500,fmt.Errorf("failed to send email after the %d attempts,error: %v",maxRetries,retryError)
}