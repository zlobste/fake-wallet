package resources

import "github.com/zlobste/fake-wallet/internal/data"

type GetUserTxsResponse struct {
	Transactions []data.Transaction `json:"transactions"`
}
