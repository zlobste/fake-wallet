package handlers

import (
	"fmt"
	"github.com/zlobste/fake-wallet/internal/app/resources"
	"github.com/zlobste/fake-wallet/internal/app/utils"
	"math/big"
	"net/http"
)

func TransferFunds(w http.ResponseWriter, r *http.Request) {
	log := utils.Log(r)

	userID, ok := r.Context().Value(utils.UserIDTag).(int64)
	if !ok {
		log.Debug(fmt.Sprintf("Failed to parse user id from context: %v", userID))
		utils.Respond(w, http.StatusInternalServerError, utils.Message("Something went wrong"))
		return
	}

	request, err := resources.NewTransferFunds(r)
	if err != nil {
		log.WithError(err).Debug(fmt.Sprintf("Invalid request to transfer funds, user id: %v", userID))
		utils.Respond(w, http.StatusBadRequest, utils.Message(err.Error()))
		return
	}

	senderWallet, err := utils.Wallets(r).GetByAddress(request.From)
	if err != nil {
		log.WithError(err).Debug(fmt.Sprintf("Failed to find sender wallet: %s", request.From))
		utils.Respond(w, http.StatusInternalServerError, utils.Message(err.Error()))
		return
	}
	if senderWallet == nil {
		log.Debug(fmt.Sprintf("Can not find sender wallet: %s", request.From))
		utils.Respond(w, http.StatusNotFound, utils.Message("Unknown sender wallet"))
		return
	}

	if senderWallet.OwnerID != userID {
		log.Debug(fmt.Sprintf("It is forbidden to use wallet: %s", senderWallet.Address))
		utils.Respond(w, http.StatusForbidden, utils.Message("Bad access"))
		return
	}

	receiverWallet, err := utils.Wallets(r).GetByAddress(request.To)
	if err != nil {
		log.WithError(err).Debug(fmt.Sprintf("Failed to find receiver wallet: %s", request.To))
		utils.Respond(w, http.StatusInternalServerError, utils.Message(err.Error()))
		return
	}
	if receiverWallet == nil {
		log.Debug(fmt.Sprintf("Can not find receiver wallet: %s", request.To))
		utils.Respond(w, http.StatusNotFound, utils.Message("Unknown receiver wallet"))
		return
	}

	if senderWallet.AssetID != receiverWallet.AssetID {
		log.Debug(fmt.Sprintf("Unsupported asset for %s", request.To))
		utils.Respond(w, http.StatusForbidden, utils.Message(fmt.Sprintf("Unsupported asset for %s", request.To)))
		return
	}

	asset, err := utils.Assets(r).GetByID(senderWallet.AssetID)
	if err != nil {
		log.WithError(err).Debug(fmt.Sprintf("Failed to find asset: %v", senderWallet.AssetID))
		utils.Respond(w, http.StatusInternalServerError, utils.Message(err.Error()))
		return
	}
	if asset == nil {
		log.Debug(fmt.Sprintf("Can not find asset: %v", senderWallet.AssetID))
		utils.Respond(w, http.StatusNotFound, utils.Message("Unsupported asset for wallet"))
		return
	}

	feePercentage := big.NewFloat(asset.Fee)
	senderBalance := senderWallet.Balance.Unwrap()
	txAmount := new(big.Float).SetInt(request.Amount.Unwrap())
	transactionFee, _ := new(big.Float).Mul(txAmount, feePercentage).Int(nil)
	if transactionFee.Cmp(new(big.Int)) != 1 {
		log.Debug(fmt.Sprintf("Transfer amount is too small: %v", request.Amount))
		utils.Respond(w, http.StatusBadRequest, utils.Message("Transfer amount is too small"))
		return
	}

	withdrawalAmount := new(big.Int).Sub(senderBalance, transactionFee)
	if withdrawalAmount.Cmp(senderBalance) == 1 {
		log.Debug(fmt.Sprintf("Insufficient funds for: %v", senderWallet.Address))
		utils.Respond(w, http.StatusForbidden, utils.Message("Insufficient funds"))
		return
	}

	if err := utils.Transactions(r).TransferFunds(request.From, request.To, request.Amount.Unwrap(), transactionFee); err != nil {
		utils.Respond(w, http.StatusInternalServerError, utils.Message("Failed to process the tx"))
		return
	}

	utils.Respond(
		w,
		http.StatusOK,
		utils.Message("Success!"),
	)
}
