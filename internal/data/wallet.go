package data

import (
	"crypto/rand"
	"encoding/hex"
)

type Wallet struct {
	Address string  `db:"address" structs:"address"`
	Balance *BigInt `db:"balance" structs:"balance"`
	AssetID int64   `db:"asset_id" structs:"asset_id"`
	OwnerID int64   `db:"owner_id" structs:"owner_id"`
}

func GenerateAddress(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
