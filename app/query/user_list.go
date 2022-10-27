package query

import (
	"context"
	"it-test/adapters/psql"
	"it-test/pkg/logs"
)

type ListUserRepository interface {
	ListUser(context.Context, *psql.Paginate, *psql.UserListFilter) ([]psql.User, int, error)
}

type ListUserHandler struct {
	repo ListUserRepository
}

type ListUser struct {
	Paginate *psql.Paginate
	Filter   *psql.UserListFilter
}

func NewListUserHandler(repo ListUserRepository) *ListUserHandler {
	return &ListUserHandler{repo: repo}
}

func (h *ListUserHandler) Handle(ctx context.Context, cmd *ListUser) (users []psql.User, count int, err error) {
	defer func() {
		logs.LogCommandExecution("ListUserHandler", cmd, err)
	}()

	return h.repo.ListUser(ctx, cmd.Paginate, cmd.Filter)
}
