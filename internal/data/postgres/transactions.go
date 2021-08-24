package postgres

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"github.com/zlobste/fake-wallet/internal/data"
	"math/big"
	"time"
)

const (
	transactionsTable = "transactions"
)

type TransactionsStorage interface {
	New() TransactionsStorage
	GetByID(id int64) (*data.Transaction, error)
	TransferFunds(from, to string, amount, fee *big.Int) error
	GetUserTransactions(userID int64, limit, offset uint64) ([]data.Transaction, error)
}

type transactionsStorage struct {
	db  *sql.DB
	sql squirrel.SelectBuilder
}

var txsSelect = squirrel.Select(all).From(transactionsTable).PlaceholderFormat(squirrel.Dollar)

func (s *transactionsStorage) New() TransactionsStorage {
	return NewTransactionsStorage(s.db)
}

func NewTransactionsStorage(db *sql.DB) TransactionsStorage {
	return &transactionsStorage{
		db:  db,
		sql: txsSelect.RunWith(db),
	}
}

func (s *transactionsStorage) Get() (*data.Transaction, error) {
	rowScanner := s.sql.QueryRow()
	model := data.Transaction{}
	err := rowScanner.Scan(
		&model.ID,
		&model.Sender,
		&model.Receiver,
		&model.Amount,
		&model.Fee,
		&model.Time,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "failed to query model")
	} else if err == sql.ErrNoRows {
		return nil, nil
	}
	return &model, nil
}

func (s *transactionsStorage) Select() ([]data.Transaction, error) {
	rows, err := s.sql.RunWith(s.db).Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var models []data.Transaction

	for rows.Next() {
		model := data.Transaction{}
		err := rows.Scan(
			&model.ID,
			&model.Sender,
			&model.Receiver,
			&model.Amount,
			&model.Fee,
			&model.Time,
		)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, nil
}

func (s *transactionsStorage) GetByID(id int64) (*data.Transaction, error) {
	s.sql = s.sql.Where(squirrel.Eq{"id": id})
	return s.Get()
}

func (s *transactionsStorage) newInsert() squirrel.InsertBuilder {
	return squirrel.Insert(transactionsTable).RunWith(s.db).PlaceholderFormat(squirrel.Dollar)
}

func (s *transactionsStorage) TransferFunds(from, to string, amount, fee *big.Int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin sql tx")
	}

	qWithdraw, args, err := squirrel.Update(walletsTable).Set(
		"balance",
		squirrel.Expr(fmt.Sprintf("balance - %s", new(big.Int).Add(amount, fee).String())),
	).Where(squirrel.Expr(fmt.Sprintf("address = '%s'", from))).ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to parse sql")
	}
	_, err = tx.Exec(qWithdraw, args...)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to withdraw funds from: %s", from))
	}

	qDeposit, args, err := squirrel.Update(walletsTable).Set(
		"balance",
		squirrel.Expr(fmt.Sprintf("balance + %s", amount.String())),
	).Where(squirrel.Expr(fmt.Sprintf("address = '%s'", to))).ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to parse sql")
	}
	_, err = tx.Exec(qDeposit, args...)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to transfer funds to: %s", to))
	}

	newTransaction := data.Transaction{
		Sender:   from,
		Receiver: to,
		Amount:   data.NewBigInt(amount),
		Fee:      data.NewBigInt(fee),
		Time:     time.Now().UTC(),
	}
	qCreateTx, args, err := s.newInsert().SetMap(structs.Map(newTransaction)).ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to parse sql")
	}
	_, err = tx.Exec(qCreateTx, args...)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to insert tx: %v", newTransaction))
	}

	return errors.Wrap(tx.Commit(), "failed to commit sql tx")
}

func (s *transactionsStorage) GetUserTransactions(userID int64, limit, offset uint64) ([]data.Transaction, error) {
	s.sql = s.sql.Where(fmt.Sprintf("sender in (select address from wallets where owner_id = %v)", userID)).Limit(limit).Offset(offset)
	return s.Select()
}
