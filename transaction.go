package queries

import "github.com/jmoiron/sqlx"

type Transaction[T QueryTypes] struct {
  Handler QueryHandler[T]
  Tx *sqlx.Tx
}

func NewTransaction[T QueryTypes](handler QueryHandler[T], tx *sqlx.Tx) *Transaction[T] {
  return &Transaction[T] { Handler : handler, Tx: tx}
}
