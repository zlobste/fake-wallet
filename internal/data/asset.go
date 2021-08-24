package data

type Asset struct {
	ID     int64   `db:"id" structs:"-"`
	Symbol string  `db:"symbol" structs:"symbol"`
	Name   string  `db:"name" structs:"name"`
	Fee    float64 `db:"fee" structs:"fee"`
}