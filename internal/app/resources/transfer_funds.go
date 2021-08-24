package resources

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/zlobste/fake-wallet/internal/data"
	"math/big"
	"net/http"
)

type TransferFundsRequest struct {
	From   string       `json:"from"`
	To     string       `json:"to"`
	Amount *data.BigInt `json:"amount"`
}

func NewTransferFunds(r *http.Request) (TransferFundsRequest, error) {
	var request TransferFundsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal transfer request")
	}
	return request, validateTransferFunds(request)
}

func validateTransferFunds(r TransferFundsRequest) error {
	if r.From == "" || r.To == "" {
		return errors.Errorf("Invalid wallet address")
	}

	if r.Amount.Unwrap().Cmp(new(big.Int)) != 1 {
		return errors.Errorf("Invalid transfer amount")
	}

	return nil
}
