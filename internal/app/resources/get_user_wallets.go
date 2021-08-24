package resources

import "github.com/zlobste/fake-wallet/internal/data"

type GetUserWalletsResponse struct {
	Wallets []data.Wallet `json:"wallets"`
}
