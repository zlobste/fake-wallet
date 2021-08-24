package postgres

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"github.com/zlobste/fake-wallet/internal/data"
	"math/big"
)

const (
	walletsTable = "wallets"
)

type WalletsStorage interface {
	New() WalletsStorage
	GetByAddress(address string) (*data.Wallet, error)
	GetUserWallets(ownerID int64) ([]data.Wallet, error)
	Create(wallet data.Wallet) error
	TransferFunds(sender, receiver string, amount *big.Int, fee *big.Int) error
}

type walletsStorage struct {
	db  *sql.DB
	sql squirrel.SelectBuilder
}

var walletsSelect = squirrel.Select(all).From(walletsTable).PlaceholderFormat(squirrel.Dollar)


func (s *walletsStorage) New() WalletsStorage {
	return NewWalletsStorage(s.db)
}

func NewWalletsStorage(db *sql.DB) WalletsStorage {
	return &walletsStorage{
		db:  db,
		sql: walletsSelect.RunWith(db),
	}
}

func (s *walletsStorage) Get() (*data.Wallet, error) {
	rowScanner := s.sql.QueryRow()
	model := data.Wallet{}
	err := rowScanner.Scan(
		&model.Address,
		&model.AssetID,
		&model.Balance,
		&model.OwnerID,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "failed to query model")
	} else if err == sql.ErrNoRows {
		return nil, nil
	}
	return &model, nil
}

func (s *walletsStorage) Select() ([]data.Wallet, error) {
	rows, err := s.sql.RunWith(s.db).Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var models []data.Wallet

	for rows.Next() {
		model := data.Wallet{}
		err := rows.Scan(
			&model.Address,
			&model.AssetID,
			&model.Balance,
			&model.OwnerID,
		)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, nil
}

func (s *walletsStorage) GetByAddress(address string) (*data.Wallet, error) {
	s.sql = s.sql.Where(squirrel.Eq{"address": address})
	return s.Get()
}

func (s *walletsStorage) Create(wallet data.Wallet) error {
	_, err := s.newInsert().SetMap(structs.Map(wallet)).Exec()
	if err != nil {
		return errors.Wrap(err, "failed to insert wallet")
	}
	return nil
}

func (s *walletsStorage) GetUserWallets(ownerID int64) ([]data.Wallet, error) {
	s.sql = s.sql.Where(squirrel.Eq{"owner_id": ownerID})
	return s.Select()
}

func (s *walletsStorage) TransferFunds(sender, receiver string, amount *big.Int, fee *big.Int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin sql tx")
	}

	withdraw := data.NewBigInt(new(big.Int).Add(amount, fee))
	qWithdraw := squirrel.Update(walletsTable).Set("balance", fmt.Sprintf("balance - %v", withdraw)).Where(squirrel.Eq{"sender": sender})
	_, err = tx.Exec(qWithdraw.ToSql())
	if err != nil {
		return errors.Wrap(err, "failed to withdraw")
	}

	qDeposit :=  squirrel.Update(walletsTable).Set("balance", fmt.Sprintf("balance + %v", amount)).Where(squirrel.Eq{"receiver": receiver})
	_, err = tx.Exec(qDeposit.ToSql())
	if err != nil {
		return errors.Wrap(err, "failed to deposit")
	}

	return errors.Wrap(tx.Commit(), "failed to commit sql tx")
}

func (s *walletsStorage) newInsert() squirrel.InsertBuilder {
	return squirrel.Insert(walletsTable).RunWith(s.db).PlaceholderFormat(squirrel.Dollar)
}