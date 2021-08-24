package handlers

import (
	"fmt"
	"github.com/zlobste/fake-wallet/internal/app/resources"
	"github.com/zlobste/fake-wallet/internal/app/utils"
	"net/http"
)

func GetUserWallets(w http.ResponseWriter, r *http.Request) {
	log := utils.Log(r)

	userID, ok := r.Context().Value(utils.UserIDTag).(int64)
	if !ok {
		log.Debug(fmt.Sprintf("Failed to parse user id from context: %v", userID))
		utils.Respond(w, http.StatusInternalServerError, utils.Message("Something went wrong"))
		return
	}

	wallets, err := utils.Wallets(r).GetUserWallets(userID)
	if err != nil {
		log.WithError(err).Debug(fmt.Sprintf("Failed to get user wallets, user id: %v", userID))
		utils.Respond(w, http.StatusInternalServerError, utils.Message("Something went wrong"))
		return
	}
	if wallets == nil {
		log.WithError(err).Debug(fmt.Sprintf("Unable to find wallets, user id: %v", userID))
		utils.Respond(w, http.StatusNotFound, utils.Message("Unable to find wallets"))
		return
	}

	utils.Respond(
		w,
		http.StatusOK,
		utils.Message(
			resources.GetUserWalletsResponse{
				Wallets: wallets,
			},
		),
	)
}
