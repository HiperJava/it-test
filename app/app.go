package app

import (
	"it-test/app/query"
	"it-test/config"
)

type Application struct {
	Commands  *Commands
	Queries   *Queries
	AppConfig *config.AppConfig
}

type Commands struct {
}

type Queries struct {
	GetUserCount *query.GetUserCountHandler
	InsertUser   *query.CreateUserHandler
	UpdateUser   *query.UpdateUserHandler
	ListUser     *query.ListUserHandler
}
