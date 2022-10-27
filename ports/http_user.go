package ports

import (
	"encoding/json"
	"errors"
	"it-test/adapters/psql"
	"it-test/app/query"
	"it-test/domain"
	"it-test/pkg/server/httperr"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

func (h HTTPServer) PostUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var createUserRequest = new(CreateUser)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&createUserRequest)
	if err != nil {
		httperr.InternalError(domain.ErrorInternalServerErrorLabel, createUser, uuid.NewString(), err, w, r)
		return
	}

	userModel := psql.User{
		UserName:  createUserRequest.UserName,
		LastName:  createUserRequest.LastName,
		FirstName: createUserRequest.FirstName,
		Password:  createUserRequest.Password,
		Email:     createUserRequest.Email,
		Mobile:    createUserRequest.Mobile,
		ASZF:      createUserRequest.Aszf,
	}

	err = h.app.Queries.InsertUser.Handle(ctx, &query.CreateUser{Model: &userModel})
	if err != nil {
		httperr.InternalError(domain.ErrorInternalServerErrorLabel, createUser, uuid.NewString(), err, w, r)
		return
	}

	response := GetUser{
		Aszf:      userModel.ASZF,
		Email:     userModel.Email,
		FirstName: userModel.FirstName,
		Id:        uuid.MustParse(userModel.ID),
		LastName:  userModel.LastName,
		Mobile:    userModel.Mobile,
		UserName:  userModel.UserName,
	}

	render.Respond(w, r, response)
}

func (h HTTPServer) GetUserList(w http.ResponseWriter, r *http.Request, params GetUserListParams) {
	ctx := r.Context()

	paginate := psql.Paginate{
		PageIndex: params.PageIndex,
		PageSize:  params.Limit,
		OrderBy:   params.OrderBy,
		Order:     params.Order,
	}
	userFilter := psql.UserListFilter{Email: params.EmailFilter}

	users, count, err := h.app.Queries.ListUser.Handle(ctx, &query.ListUser{Paginate: &paginate, Filter: &userFilter})
	if err != nil {
		httperr.InternalError(domain.ErrorInternalServerErrorLabel, listUser, uuid.NewString(), err, w, r)
		return
	}

	userListResults := make([]UserListItem, 0, len(users))

	for _, userModel := range users {
		userListItem := UserListItem{
			Aszf:      userModel.ASZF,
			Email:     userModel.Email,
			FirstName: userModel.FirstName,
			Id:        uuid.MustParse(userModel.ID),
			LastName:  userModel.LastName,
			Mobile:    userModel.Mobile,
			UserName:  userModel.UserName,
		}

		userListResults = append(userListResults, userListItem)
	}

	response := UserList{
		Results:       userListResults,
		ResultsLength: count,
	}

	render.Respond(w, r, response)
}

func (h HTTPServer) UpdateUserDetails(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	var createUserRequest = new(UpdateUser)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&createUserRequest)
	if err != nil {
		httperr.InternalError(domain.ErrorInternalServerErrorLabel, updateUser, uuid.NewString(), err, w, r)
		return
	}

	userModel := psql.User{
		ID:        id,
		UserName:  createUserRequest.UserName,
		LastName:  createUserRequest.LastName,
		FirstName: createUserRequest.FirstName,
		Password:  createUserRequest.Password,
		Mobile:    createUserRequest.Mobile,
	}

	err = h.app.Queries.UpdateUser.Handle(ctx, &query.UpdateUser{Model: &userModel})
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			errorBody := httperr.NewErrorMessageBody(domain.ErrorUserNotFound, "server", updateUser, uuid.NewString())
			httperr.NotFound(*errorBody, err, w, r)
			return
		}

		httperr.InternalError(domain.ErrorInternalServerErrorLabel, updateUser, uuid.NewString(), err, w, r)
		return
	}

	response := GetUser{
		Aszf:      userModel.ASZF,
		Email:     userModel.Email,
		FirstName: userModel.FirstName,
		Id:        uuid.MustParse(userModel.ID),
		LastName:  userModel.LastName,
		Mobile:    userModel.Mobile,
		UserName:  userModel.UserName,
	}

	render.Respond(w, r, response)
}

func (h HTTPServer) Count(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, err := h.app.Queries.GetUserCount.Handle(ctx, &query.GetUserCount{})
	if err != nil {
		httperr.InternalError(domain.ErrorInternalServerErrorLabel, getUserCount, uuid.NewString(), err, w, r)
		return
	}
	response := Count{Count: &result}

	render.Respond(w, r, response)
}
