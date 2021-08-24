package data

import "time"

type Transaction struct {
	ID       int64     `db:"id" structs:"-"`
	Sender   string    `db:"sender" structs:"sender"`
	Receiver string    `db:"receiver" structs:"receiver"`
	Amount   *BigInt   `db:"amount" structs:"amount"`
	Fee      *BigInt   `db:"fee" structs:"fee"`
	Time     time.Time `db:"time" structs:"time"`
}
