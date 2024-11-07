package sqlx_s

import (
	queries "github.com/nam2184/generic-queries"
  util "github.com/nam2184/generic-queries/utils"
  "fmt"

)

func Insert[T queries.QueryTypes](getRow bool) queries.QueryHandlerFunc[T] {
    return func(q *queries.Query[T]) error {
        if q.A != nil {
            // Handle slice of T
            if q.Tx != nil {
                for index, item := range q.A {
                    
                    fields, _ := util.Fields[T](item)
                    args, err := util.GetArgs[T](item); if err != nil {
                      return err
                    }
                     
                    if getRow == true {
                      placeholders, _ := util.GeneratePositionalParams[T](item)
                      query := fmt.Sprintf("INSERT INTO %s (%s) VALUES( %s ) RETURNING *", 
                                        item.TableName(), 
                                        fields,
                                        placeholders, 
                                        )
                      var row T
                      err = q.Tx.QueryRowx(query, args...).StructScan(&row); if err != nil {
                          return err
                      }
                      if util.IsZero[T](row) {
                        return fmt.Errorf("No row returned for insert %d", index)
                      }
                      q.Rows = append(q.Rows, row)
                    } else {
                      placeholders, _ := util.GenerateNamedParams[T](item)
                      query := fmt.Sprintf("INSERT INTO %s (%s) VALUES( %s )", 
                                        item.TableName(), 
                                        fields,
                                        placeholders, 
                                        )
 
                      _, err = q.Tx.NamedExec(query, &item)
                    } 
                    if err != nil {
                        if rollbackErr := q.Tx.Rollback(); rollbackErr != nil {
                            return fmt.Errorf("Failed to rollback transaction: %s", rollbackErr)
                        }
                        return err
                    }
                }
            } else {
                return fmt.Errorf("Should not handle slice without transaction") 
            }
        }
        if len(q.Rows) == 0 { 
          return fmt.Errorf("Error inserting rows, no rows returned as a result")
        }
        return nil
    }
}

func Select[T queries.QueryTypes](constraint string) queries.QueryHandlerFunc[T] {
    return func(q *queries.Query[T]) error {
        // Handle slice of T
        if q.Tx != nil {
            var item T
            query := fmt.Sprintf("SELECT * FROM %s WHERE %s", 
                                    item.TableName(),
                                    constraint, 
                                )
            
            err := q.Tx.Select(&q.Rows, query)
            if err != nil {
                if rollbackErr := q.Tx.Rollback(); rollbackErr != nil {
                    return fmt.Errorf("Failed to rollback transaction: %s", rollbackErr)
                }
                return err
            }
        } else {
                return fmt.Errorf("Should not handle slice without transaction") 
            }
        return nil
        }
}

func SelectOffset[T queries.QueryTypes](args map[string]interface{}, limit int, skip int, sort_by string, order string, constraint Constraint) queries.QueryHandlerFunc[T] {
    return func(q *queries.Query[T]) error {
        if q.Tx != nil {
            var item T
            var err error
            q.Rows = make([]T, 0, len(q.A))
            number, err := constraint.GetFinalPlaceholder(); if err != nil {
                  number = 1
            }
            if len(args) == 0 {
              query := fmt.Sprintf(
                  "SELECT * FROM %s WHERE %s ORDER BY $%d %s LIMIT $%d OFFSET $%d",
                  item.TableName(),
                  constraint.constraint,
                  number + 1,
                  order,
                  number + 2,
                  number + 3,
              )
              full_args := append(constraint.values, sort_by, limit, skip)
              err = q.Tx.Select(&q.Rows, query, full_args...)
            } else {
              filters , values := util.GenerateFilterString(args, number, limit, skip, sort_by)
              query := fmt.Sprintf(
                  "SELECT * FROM %s WHERE %s AND %s ORDER BY $%d %s LIMIT $%d OFFSET $%d",
                  item.TableName(),
                  filters,
                  constraint.constraint,
                  number + 1,
                  order,
                  number + 2,
                  number + 3,
              )

              full_args := append(constraint.values, sort_by, limit, skip, values)
              err = q.Tx.Select(&q.A, query, full_args...)
            }
            if err != nil {
                if rollbackErr := q.Tx.Rollback(); rollbackErr != nil {
                    return fmt.Errorf("Failed to rollback transaction: %s", rollbackErr)
                }
                return err
            }
        } else {
                return fmt.Errorf("Should not handle slice without transaction") 
            }
        return nil
        }
}
func SelectOffset2[T queries.QueryTypes](args map[string]interface{}, limit int, skip int, sort_by string, order string, constraint string) queries.QueryHandlerFunc[T] {
    return func(q *queries.Query[T]) error {
        if q.Tx != nil {
            var item T
            var err error
            q.Rows = make([]T, 0, len(q.A))
            if len(args) == 0 {
              query := fmt.Sprintf(
                  "SELECT * FROM %s WHERE %s ORDER BY $1 %s LIMIT $2 OFFSET $3",
                  item.TableName(),
                  constraint,
                  order,
              )
              err = q.Tx.Select(&q.Rows, query, sort_by, limit, skip)
            } else {
              filters , values := util.GenerateFilterString(args, limit, skip, sort_by)
              query := fmt.Sprintf(
                  "SELECT * FROM %s WHERE %s AND %s ORDER BY $1 %s LIMIT $2 OFFSET $3",
                  item.TableName(),
                  filters,
                  constraint,
                  order,
              )
              err = q.Tx.Select(&q.A, query, values...)
            }
            if err != nil {
                if rollbackErr := q.Tx.Rollback(); rollbackErr != nil {
                    return fmt.Errorf("Failed to rollback transaction: %s", rollbackErr)
                }
                return err
            }
        } else {
                return fmt.Errorf("Should not handle slice without transaction") 
            }
        return nil
        }
}


func Update[T queries.QueryTypes](constraint string) queries.QueryHandlerFunc[T] {
    return func(q *queries.Query[T]) error {
        if q.A == nil {
            return fmt.Errorf("No items provided to update")
        }
        
        if q.Tx == nil {
            return fmt.Errorf("Update operation on slice should have an active transaction")
        }
        
        q.Rows = make([]T, 0, len(q.A))

        for index, item := range q.A {
            fields, _ := util.FieldsAndPlaceholders[T](item)
            args, err := util.GetArgs[T](item); if err != nil {
              return err
            }
            query := fmt.Sprintf("UPDATE %s SET %s WHERE %s RETURNING *", 
                                 item.TableName(),
                                 fields,
                                 constraint)
            var updatedRow T
            err = q.Tx.QueryRowx(query, args...).StructScan(&updatedRow); if err != nil {
                if rollbackErr := q.Tx.Rollback(); rollbackErr != nil {
                    return fmt.Errorf("Failed to rollback transaction: %s", rollbackErr)
                }
                return err
            }
            if util.IsZero[T](updatedRow) {
              return fmt.Errorf("No row returned for %d update", index)
            }
            q.Rows = append(q.Rows, updatedRow)
        }

        return nil
    }
}

