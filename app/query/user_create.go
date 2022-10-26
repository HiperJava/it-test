package query

import (
	"context"
	"it-test/adapters/psql"
	"it-test/pkg/logs"
)

type CreateUserRepository interface {
	InsertUser(ctx context.Context, user *psql.User) error
}

type CreateUserHandler struct {
	repo CreateUserRepository
}

type CreateUser struct {
	Model *psql.User
}

func NewCreateUserHandler(repo CreateUserRepository) *CreateUserHandler {
	return &CreateUserHandler{repo: repo}
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd *CreateUser) (err error) {
	defer func() {
		logs.LogCommandExecution("CreateUserHandler", cmd, err)
	}()

	return h.repo.InsertUser(ctx, cmd.Model)
}
