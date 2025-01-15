package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
)


type Role struct{
	Id int64 `json:"id"`
	Name string `json:"name"`
	Level int64 `json:"level"`
	Description string `json:"description"`
}


type RoleStore struct{
	db *sql.DB
}

func (s *RoleStore) GetByName(ctx context.Context,name string) (*Role,error){
	query:=`
	SELECT id,name,level,description FROM roles where name=$1
	`
	var role Role
	ctx,cancel:=context.WithTimeout(ctx,QueryTimeoutDuration)
	defer cancel()
	err:=s.db.QueryRowContext(ctx,query,role).Scan(&role.Id,&role.Name,&role.Level,&role.Description)
	if err!=nil{
		switch {
		case	errors.Is(err,sql.ErrNoRows):
			return nil,ErrorNotFound
		default:
			return nil,err
		}
	}
	log.Printf("post: %v",role)
	return &role,err

} 