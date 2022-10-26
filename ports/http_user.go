package ports

import (
	"encoding/json"
	"fmt"
	"it-test/adapters/psql"
	"it-test/app/query"
	"it-test/domain"
	"it-test/pkg/server/httperr"
	"net/http"

	"github.com/go-chi/render"
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

	fmt.Printf("Request %+v\n\n", createUserRequest)

	userModel := psql.User{
		UserName:  createUserRequest.UserName,
		LastName:  createUserRequest.LastName,
		FirstName: createUserRequest.FirstName,
		Password:  createUserRequest.Password,
		Email:     createUserRequest.Email,
		Mobile:    createUserRequest.Mobile,
		ASZF:      createUserRequest.Aszf,
	}

	fmt.Printf("Model %+v\n\n", userModel)

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
	//TODO implement me
	panic("implement me")
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
