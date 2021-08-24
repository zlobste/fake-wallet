package handlers

import (
	"fmt"
	"github.com/zlobste/fake-wallet/internal/app/resources"
	"github.com/zlobste/fake-wallet/internal/app/utils"
	"github.com/zlobste/fake-wallet/internal/data"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	log := utils.Log(r)

	request, err := resources.NewRegister(r)
	if err != nil {
		log.WithError(err).Debug("Failed to parse login request")
		utils.Respond(w, http.StatusBadRequest, utils.Message(err.Error()))
		return
	}

	users := utils.Users(r)
	user, err := users.GetByEmail(request.Email)
	if err != nil {
		log.WithError(err).Debug("Failed to get user by email")
		utils.Respond(
			w,
			http.StatusInternalServerError,
			utils.Message(fmt.Sprintf("Something bad happened, error: %s", err.Error())),
		)
		return
	}
	if user != nil {
		log.WithError(err).Debug("Specified user exist already")
		utils.Respond(w, http.StatusNotFound, utils.Message("Specified user exist already"))
		return
	}

	newUser := data.User{
		Name:     request.Name,
		Surname:  request.Surname,
		Email:    request.Email,
		Password: request.Password,
	}

	if err := newUser.EncryptPassword(); err != nil {
		log.WithError(err).Debug("Failed to encrypt password")
		utils.Respond(w, http.StatusNotFound, utils.Message(err.Error()))
		return
	}

	if err := users.Create(newUser); err != nil {
		log.WithError(err).Debug("Failed to create user")
		utils.Respond(w, http.StatusInternalServerError, utils.Message(err.Error()))
		return
	}

	utils.Respond(
		w,
		http.StatusOK,
		utils.Message("Success!"),
	)
}
