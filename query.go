package queries

import (
	"github.com/jmoiron/sqlx"
)

type Query[T QueryTypes] struct {
    A           []T
    Rows        []T
    Total       int
    Tx          *sqlx.Tx
    Q           []string
}

func NewQueries[T QueryTypes](a []T, tx *sqlx.Tx) *Query[T] {
    return &Query[T]{
        A:  a,
        Tx: tx,
    }
}

type QueryTypes interface {
  TableName() string
  Id()        interface{}
}

type QueryHandlerFunc[T QueryTypes] func(*Query[T]) error


func (f QueryHandlerFunc[T]) HandleQuery(q *Query[T]) error {
    err := f(q)
    return err
}

type QueryHandler[T QueryTypes] interface {
    HandleQuery(q *Query[T]) error
}
