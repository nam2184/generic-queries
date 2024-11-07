package sqlx_s

import (
	queries "github.com/nam2184/generic-queries"

	"github.com/jmoiron/sqlx"
)

/*
These generic queries create a transaction struct if none is passed, thus defining the query for you,
suitable for one specific query done routes and routes that handle multiple queries.QueryTypes
*/

//If you want to save memory by having specific *queries.Transactions passed, define InsertQuery[T](nil, tran, data)


func InsertQuery[T queries.QueryTypes](tx *sqlx.Tx, tran *queries.Transaction[T], data []T) (*queries.Query[T], error) {
    if tran == nil {
        // Create a new transaction if not provided
        tran = queries.NewTransaction[T](Insert[T](false), tx)
    }
    
    qs := queries.NewQueryMany[T](data, tran.Tx)
    err := tran.Handler.HandleQuery(qs)
    return qs, err
}

func InsertQueryRow[T queries.QueryTypes](tx *sqlx.Tx, tran *queries.Transaction[T], data []T) (*queries.Query[T], error) {
    if tran == nil {
        // Create a new transaction if not provided
        tran = queries.NewTransaction[T](Insert[T](true), tx)
    }
    
    qs := queries.NewQueryMany[T](data, tran.Tx)
    err := tran.Handler.HandleQuery(qs)   
 
    return qs, err
}


func SelectQuery[T queries.QueryTypes](tx *sqlx.Tx, constraint string, tran *queries.Transaction[T], data []T) (*queries.Query[T], error) {
    if tran == nil {
        // Create a new transaction if not provided
        tran = queries.NewTransaction[T](Select[T](constraint), tx)
    }
    
    qs := queries.NewQueryMany[T](data, tran.Tx)
    err := tran.Handler.HandleQuery(qs) 
    return qs, err
}

func SelectOffsetQuery[T queries.QueryTypes](tx *sqlx.Tx, 
                                              limit int, 
                                              skip int, 
                                              sort_by string, 
                                              order string, 
                                              args map[string]interface{}, 
                                              constraint Constraint, 
                                              tran *queries.Transaction[T], 
                                              data []T) (*queries.Query[T], error) {
    if tran == nil {
        // Create a new transaction if not provided
        tran = queries.NewTransaction[T](SelectOffset[T](args, limit, skip, sort_by, order, constraint), tx)
    }
     
    qs := queries.NewQueryMany[T](data, tran.Tx)
    err := tran.Handler.HandleQuery(qs) 
    return qs, err
}


func UpdateQuery[T queries.QueryTypes](tx *sqlx.Tx, constraint string, tran *queries.Transaction[T], data []T) (*queries.Query[T], error) {
    if tran == nil {
        tran = queries.NewTransaction[T](Update[T](constraint), tx)
    }

    qs := queries.NewQueryMany[T](data, tran.Tx)
    err := tran.Handler.HandleQuery(qs) 
    return qs, err
}


