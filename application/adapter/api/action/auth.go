package action

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"github.com/jeferagudeloc/gorilla-sessions-auth/application/adapter/api/response"
	"github.com/jeferagudeloc/gorilla-sessions-auth/application/adapter/logger"
	"github.com/jeferagudeloc/gorilla-sessions-auth/application/adapter/logging"
	"github.com/jeferagudeloc/gorilla-sessions-auth/application/usecase"
	"github.com/jeferagudeloc/gorilla-sessions-auth/domain"
)

type AuthAction struct {
	uc          usecase.AuthUseCase
	log         logger.Logger
	cookieStore *sessions.CookieStore
}

func NewAuthAction(uc usecase.AuthUseCase, log logger.Logger, cookieStore *sessions.CookieStore) AuthAction {
	return AuthAction{
		uc:          uc,
		log:         log,
		cookieStore: cookieStore,
	}
}

func (fova AuthAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "login"

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

	var requestBody *domain.Auth
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

	output, err := fova.uc.Execute(r.Context(), domain.Auth{
		Username: requestBody.Username,
		Password: requestBody.Password,
	})

	if err != nil {
		switch err.Error() {

		case domain.ErrPasswordDoesNotMatch:
			logging.NewError(
				fova.log,
				err,
				logKey,
				http.StatusUnauthorized,
			).Log(domain.ErrPasswordDoesNotMatch)
			response.NewError(err, http.StatusUnauthorized).Send(w)
			return
		case domain.ErrUserNotFound:
			logging.NewError(
				fova.log,
				err,
				logKey,
				http.StatusUnauthorized,
			).Log(domain.ErrUserNotFound)
			response.NewError(err, http.StatusUnauthorized).Send(w)
			return
		default:
			logging.NewError(
				fova.log,
				err,
				logKey,
				http.StatusInternalServerError,
			).Log("error login")

			response.NewError(err, http.StatusInternalServerError).Send(w)
			return
		}
	}

	fova.createSession(r, w, output)

	logging.NewInfo(fova.log, logKey, http.StatusOK).Log("login success")
	response.NewSuccess(output, http.StatusOK).Send(w)
}

func (fova AuthAction) createSession(r *http.Request, w http.ResponseWriter, user *domain.User) bool {
	claims := jwt.MapClaims{
		"name":     user.Name,
		"lastname": user.LastName,
		"email":    user.Email,
		"status":   user.Status,
		"profile": jwt.MapClaims{
			"name":        user.Profile.Name,
			"permissions": user.Profile.Permissions,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(os.Getenv("SECRET_KEY_TOKEN")))

	session, err := fova.cookieStore.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: false,
	}
	session.ID = user.Name
	session.Values["token"] = signedToken

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}
