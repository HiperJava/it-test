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
	//TODO implement me
	panic("implement me")
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
