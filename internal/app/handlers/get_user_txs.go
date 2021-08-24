package handlers

import (
	"fmt"
	"github.com/zlobste/fake-wallet/internal/app/resources"
	"github.com/zlobste/fake-wallet/internal/app/utils"
	"net/http"
	"strconv"
)

const (
	DefaultRawsLimit = 10
	DefaultOffset    = 0
)

func GetUserTxs(w http.ResponseWriter, r *http.Request) {
	log := utils.Log(r)

	userID, ok := r.Context().Value(utils.UserIDTag).(int64)
	if !ok {
		log.Debug(fmt.Sprintf("Failed to parse user id from context: %v", userID))
		utils.Respond(w, http.StatusInternalServerError, utils.Message("Something went wrong"))
		return
	}

	var offset, limit uint64 = DefaultOffset, DefaultRawsLimit
	offsetRaw, limitRaw := r.URL.Query().Get("offset"), r.URL.Query().Get("limit")
	if offsetRaw != "" {
		var err error
		offset, err = strconv.ParseUint(offsetRaw, 10, 64)
		if err != nil {
			log.WithField("offset", offsetRaw).WithError(err).Debug("Bad offset was provided")
			utils.Respond(w, http.StatusBadRequest, utils.Message("Bad offset value"))
			return
		}
	}
	if limitRaw != "" {
		var err error
		limit, err = strconv.ParseUint(limitRaw, 10, 64)
		if err != nil {
			log.WithField("count", limitRaw).WithError(err).Debug("Bad count was provided")
			utils.Respond(w, http.StatusBadRequest, utils.Message("Bad count value"))
			return
		}
	}

	transactions, err := utils.Transactions(r).GetUserTransactions(userID, limit, offset)
	if err != nil {
		log.WithError(err).Debug(fmt.Sprintf("Failed to get user txs, user id: %v", userID))
		utils.Respond(w, http.StatusInternalServerError, utils.Message("Something went wrong"))
		return
	}
	if transactions == nil {
		log.WithError(err).Debug(fmt.Sprintf("Unable to find txs, user id: %v", userID))
		utils.Respond(w, http.StatusNotFound, utils.Message("Unable to find transactions"))
		return
	}

	utils.Respond(
		w,
		http.StatusOK,
		utils.Message(
			resources.GetUserTxsResponse{
				Transactions: transactions,
			},
		),
	)

}
