package handlers

import (
	"fmt"
	"github.com/zlobste/fake-wallet/internal/app/resources"
	"github.com/zlobste/fake-wallet/internal/app/utils"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	log := utils.Log(r)

	request, err := resources.NewLogin(r)
	if err != nil {
		log.WithError(err).Debug("failed to parse login request")
		utils.Respond(w, http.StatusBadRequest, utils.Message(err.Error()))
		return
	}

	user, err := utils.Users(r).GetByEmail(request.Email)
	if err != nil {
		log.WithError(err).Debug("Failed to get user by email")
		utils.Respond(w, http.StatusInternalServerError, utils.Message(err.Error()))
		return
	}
	if user == nil {
		log.WithError(err).Debug(fmt.Sprintf("failed to find a user with email: %s", request.Email))
		utils.Respond(w, http.StatusUnauthorized, utils.Message("Invalid access"))
		return
	}

	if !user.ComparePassword(request.Password) {
		log.WithError(err).Debug(fmt.Sprintf(fmt.Sprintf("Invalid password for user id: %v", user.ID)))
		utils.Respond(w, http.StatusUnauthorized, utils.Message("Invalid access"))
		return
	}

	token, err := utils.JWT(r).CreateJWT(user.ID)
	if err != nil {
		log.WithError(err).Debug(fmt.Sprintf(fmt.Sprintf("Failed to create jwt for user id: %v", user.ID)))
		utils.Respond(w, http.StatusInternalServerError, utils.Message("Something went wrong"))
		return
	}

	utils.Respond(
		w,
		http.StatusOK,
		utils.Message(
			resources.LoginResponse{
				Token: token,
			},
		),
	)
}
