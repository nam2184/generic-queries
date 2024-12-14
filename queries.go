package queries

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

/*
These generic create a transaction struct if none is passed, thus defining the query for you,
suitable for one specific query done routes and routes that handle multiple QueryTypes
*/

//If you want to save memory by having specific *Transactions passed, define InsertQuery[T](nil, tran, data)


func InsertQuery[T QueryTypes](tx *sqlx.Tx, tran *Transaction[T], data ...T) (*Query[T], error) {
    if tran == nil {
        // Create a new transaction if not provided
        tran = NewTransaction[T](Insert[T](false), tx)
    }
    qs := NewQueries[T](data, tran.Tx)
    err := tran.Handler.HandleQuery(qs)
    return qs, err
}

func InsertQueryRow[T QueryTypes](tx *sqlx.Tx, tran *Transaction[T], data ...T) (*Query[T], error) {
    if tran == nil {
        // Create a new transaction if not provided
        tran = NewTransaction[T](Insert[T](true), tx)
    }
    
    qs := NewQueries[T](data, tran.Tx)
    err := tran.Handler.HandleQuery(qs)   
 
    return qs, err
}

func InsertQueryOne[T QueryTypes](tx *sqlx.Tx, tran *Transaction[T], query *Query[T], data T) (error) {
    if tran == nil {
        // Create a new transaction if not provided
        return fmt.Errorf("Use a transaction, wasting memory")
    }
   
    tran.Handler = InsertOne(data)
    err := tran.Handler.HandleQuery(query)   
 
    return err
}



func SelectQuery[T QueryTypes](tx *sqlx.Tx, constraint Constraint, tran *Transaction[T]) (*Query[T], error) {
    if tran == nil {
        // Create a new transaction if not provided
        tran = NewTransaction[T](Select[T](constraint), tx)
    }
    
    qs := NewQueries[T](nil, tran.Tx)
    err := tran.Handler.HandleQuery(qs) 
    return qs, err
}

func SelectOffsetQuery[T QueryTypes](
                        tx *sqlx.Tx, 
                        limit int, 
                        skip int, 
                        sort_by string, 
                        order string, 
                        args map[string]interface{}, 
                        constraint Constraint, 
                        tran *Transaction[T]) (*Query[T], error) {
    if tran == nil {
        // Create a new transaction if not provided
        tran = NewTransaction[T](SelectOffset[T](args, limit, skip, sort_by, order, constraint), tx)
    }
     
    qs := NewQueries[T](nil, tran.Tx)
    err := tran.Handler.HandleQuery(qs) 
    return qs, err
}


func UpdateQuery[T QueryTypes](tx *sqlx.Tx, constraint Constraint, tran *Transaction[T], data ...T) (*Query[T], error) {
    if tran == nil {
        tran = NewTransaction[T](Update[T](constraint), tx)
    }

    qs := NewQueries[T](data, tran.Tx)
    err := tran.Handler.HandleQuery(qs) 
    return qs, err
}

func UpdateQueryNoConstraint[T QueryTypes](tx *sqlx.Tx, tran *Transaction[T], data ...T) (*Query[T], error) {
    if tran == nil {
        tran = NewTransaction[T](UpdateNoConstraint[T](), tx)
    }

    qs := NewQueries[T](data, tran.Tx)
    err := tran.Handler.HandleQuery(qs) 
    return qs, err
}

