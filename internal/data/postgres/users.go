package postgres

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"github.com/zlobste/fake-wallet/internal/data"
	"math/big"
)

const (
	usersTable = "users"
	all        = "*"
)

const (
	InitialBalance       = 100
	DefaultAddressLength = 32
)

var DefaultBalance = data.NewBigInt(big.NewInt(InitialBalance))

type UsersStorage interface {
	New() UsersStorage
	GetByEmail(email string) (*data.User, error)
	Create(user data.User) error
}

type usersStorage struct {
	db  *sql.DB
	sql squirrel.SelectBuilder
}

var usersSelect = squirrel.Select(all).From(usersTable).PlaceholderFormat(squirrel.Dollar)

func (s *usersStorage) New() UsersStorage {
	return NewUsersStorage(s.db)
}

func NewUsersStorage(db *sql.DB) UsersStorage {
	return &usersStorage{
		db:  db,
		sql: usersSelect.RunWith(db),
	}
}

func (s *usersStorage) Get() (*data.User, error) {
	rowScanner := s.sql.QueryRow()
	model := data.User{}
	err := rowScanner.Scan(
		&model.ID,
		&model.Name,
		&model.Surname,
		&model.Email,
		&model.Password,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "failed to query model")
	} else if err == sql.ErrNoRows {
		return nil, nil
	}
	return &model, nil
}

func (s *usersStorage) GetUserById(id uint64) (*data.User, error) {
	s.sql = s.sql.Where(squirrel.Eq{"id": id})
	return s.Get()
}

func (s *usersStorage) GetByEmail(email string) (*data.User, error) {
	s.sql = s.sql.Where(squirrel.Eq{"email": email})
	return s.Get()
}

func (s *usersStorage) newInsert() squirrel.InsertBuilder {
	return squirrel.Insert(usersTable).RunWith(s.db).PlaceholderFormat(squirrel.Dollar)
}

func (s *usersStorage) Create(user data.User) error {
	assetsStorage := NewAssetsStorage(s.db)
	assets, err := assetsStorage.Select()
	if err != nil {
		return errors.Wrap(err, "failed to get assets")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin sql tx")
	}

	qCreateUser, args, err := s.newInsert().SetMap(structs.Map(user)).Suffix("returning id").ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to parse sql")
	}
	_, err = tx.Exec(qCreateUser, args...)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	var userID int64
	qSelectUserID, args, err := squirrel.Select("id").From(usersTable).PlaceholderFormat(squirrel.Dollar).Where(squirrel.Eq{"email": user.Email}).ToSql()
	row := tx.QueryRow(qSelectUserID, args...)
	if err := row.Scan(&userID); err != nil {
		return errors.Wrap(err, "failed to get user id")
	}

	walletsInsert := squirrel.Insert(walletsTable).PlaceholderFormat(squirrel.Dollar).Columns(
		"address",
		"balance",
		"asset_id",
		"owner_id",
	)
	for _, asset := range assets {
		address := data.GenerateAddress(DefaultAddressLength)
		walletsInsert = walletsInsert.Values(
			address,
			DefaultBalance,
			asset.ID,
			userID,
		)
	}

	qCreateWallets, args, err := walletsInsert.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to parse sql")
	}

	_, err = tx.Exec(qCreateWallets, args...)
	if err != nil {
		return errors.Wrap(err, "failed to create user wallets")
	}

	return errors.Wrap(tx.Commit(), "failed to commit sql tx")
}
