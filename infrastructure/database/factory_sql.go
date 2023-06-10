package database

import (
	"errors"

	"github.com/jeferagudeloc/gorilla-sessions-auth/application/adapter/repository"
)

var (
	errInvalidSQLDatabaseInstance = errors.New("invalid sql db instance")
)

const (
	InstanceMysql int = iota
)

func NewDatabaseSQLFactory(instance int) (repository.SQL, error) {
	switch instance {
	case InstanceMysql:
		return NewPostgresqlHandler(newConfigPostgresql())
	default:
		return nil, errInvalidSQLDatabaseInstance
	}
}
