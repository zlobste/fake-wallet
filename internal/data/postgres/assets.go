package postgres

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/zlobste/fake-wallet/internal/data"
)

const (
	assetsTable = "assets"
)

var assetsSelect = squirrel.Select(all).From(assetsTable).PlaceholderFormat(squirrel.Dollar)

type AssetsStorage interface {
	New() AssetsStorage
	GetByID(id int64) (*data.Asset, error)
	Select() ([]data.Asset, error)
}

type assetsStorage struct {
	db  *sql.DB
	sql squirrel.SelectBuilder
}

func (s *assetsStorage) New() AssetsStorage {
	return NewAssetsStorage(s.db)
}

func NewAssetsStorage(db *sql.DB) AssetsStorage {
	return &assetsStorage{
		db:  db,
		sql: assetsSelect.RunWith(db),
	}
}

func (s *assetsStorage) Get() (*data.Asset, error) {
	rowScanner := s.sql.QueryRow()
	model := data.Asset{}
	err := rowScanner.Scan(
		&model.ID,
		&model.Symbol,
		&model.Name,
		&model.Fee,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "failed to query model")
	} else if err == sql.ErrNoRows {
		return nil, nil
	}
	return &model, nil
}

func (s *assetsStorage) Select() ([]data.Asset, error) {
	rows, err := s.sql.RunWith(s.db).Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var models []data.Asset

	for rows.Next() {
		model := data.Asset{}
		err := rows.Scan(
			&model.ID,
			&model.Symbol,
			&model.Name,
			&model.Fee,
		)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, nil
}

func (s *assetsStorage) GetByID(id int64) (*data.Asset, error) {
	s.sql = s.sql.Where(squirrel.Eq{"id": id})
	return s.Get()
}