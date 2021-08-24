package data

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
	"math/big"
)

// BigInt is a wrapper for saving to db.
type BigInt big.Int

// Unwrap returns underlying big.Int
func (b *BigInt) Unwrap() *big.Int {
	if b == nil {
		return nil
	}

	return (*big.Int)(b)
}

// NewBigInt.
func NewBigInt(b *big.Int) *BigInt {
	return (*BigInt)(b)
}

// Value implements driver.Valuer.
func (b BigInt) Value() (driver.Value, error) {
	return []byte(b.Unwrap().String()), nil
}

// Scan implements sql.Scanner interface.
func (b *BigInt) Scan(src interface{}) error {
	switch t := src.(type) {
	case []byte:
		z, ok := big.NewInt(0).SetString(string(t), 10)
		if !ok {
			return errors.Errorf("invalid big.Int %s", string(t))
		}

		*b = BigInt(*z)
	case string:
		z, ok := big.NewInt(0).SetString(t, 10)
		if !ok {
			return errors.Errorf("invalid big.Int %s", t)
		}

		*b = BigInt(*z)
	default:
		return errors.Errorf("unexpected input type %T", src)
	}

	return nil
}

func (b *BigInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Unwrap().String())
}

func (b *BigInt) UnmarshalJSON(input []byte) error {
	z, ok := big.NewInt(0).SetString(string(input), 10)
	if !ok {
		return errors.Errorf("invalid big.Int %s", string(input))
	}
	*b = BigInt(*z)

	return nil
}
