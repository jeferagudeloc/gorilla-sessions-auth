package action

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jeferagudeloc/gorilla-sessions-auth/application/adapter/api/response"
	"github.com/jeferagudeloc/gorilla-sessions-auth/application/adapter/logger"
	"github.com/jeferagudeloc/gorilla-sessions-auth/application/adapter/logging"
	"github.com/jeferagudeloc/gorilla-sessions-auth/application/usecase"
	"github.com/jeferagudeloc/gorilla-sessions-auth/domain"
)

type CreateUserAction struct {
	uc          usecase.CreateUserUseCase
	log         logger.Logger
	cookieStore *sessions.CookieStore
}

func NewCreateUserAction(uc usecase.CreateUserUseCase, log logger.Logger, cookieStore *sessions.CookieStore) CreateUserAction {
	return CreateUserAction{
		uc:          uc,
		log:         log,
		cookieStore: cookieStore,
	}
}

func (fova CreateUserAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "creating user"

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logging.NewError(
			fova.log,
			err,
			logKey,
			http.StatusInternalServerError,
		).Log("error reading request body")

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}

	var requestBody *domain.CreateUserRequest
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		logging.NewError(
			fova.log,
			err,
			logKey,
			http.StatusBadRequest,
		).Log("error decoding request body")

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}

	output, err := fova.uc.Execute(r.Context(), requestBody)

	if err != nil {
		switch err.Error() {
		case domain.ErrConflictUser:
			logging.NewError(
				fova.log,
				err,
				logKey,
				http.StatusInternalServerError,
			).Log(domain.ErrConflictUser)
			response.NewError(err, http.StatusConflict).Send(w)
			return
		default:
			logging.NewError(
				fova.log,
				err,
				logKey,
				http.StatusInternalServerError,
			).Log("error creating user")
			response.NewError(err, http.StatusInternalServerError).Send(w)
			return
		}
	}

	logging.NewInfo(fova.log, logKey, http.StatusOK).Log("creating user success")
	response.NewSuccess(output, http.StatusOK).Send(w)
}
