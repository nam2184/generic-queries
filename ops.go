package queries

import (
  util "github.com/nam2184/generic-queries/utils"
  "fmt"

)

func Insert[T QueryTypes](getRow bool) QueryHandlerFunc[T] {
    return func(q *Query[T]) error {
        if q.A != nil {
            // Handle slice of T
            if q.Tx != nil {
                for index, item := range q.A {
                    
                    fields, placeholders, args, err := util.GetSQLParts[T](item); if err != nil {
                      return err
                    }
                     
                    if getRow == true {
                      query := fmt.Sprintf("INSERT INTO %s (%s) VALUES( %s ) RETURNING *", 
                                        item.TableName(), 
                                        fields,
                                        placeholders, 
                                        )
                      var row T
                      
                      rowResult := q.Tx.QueryRowx(query, args...)
                      if rowResult.Err() != nil {
                          return fmt.Errorf("query execution error: %w", rowResult.Err())
                      }

                      err = rowResult.StructScan(&row)
                      if err != nil {
                          return fmt.Errorf("struct scan error: %w", err)
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

func Select[T QueryTypes](constraint Constraint) QueryHandlerFunc[T] {
    return func(q *Query[T]) error {
        // Handle slice of T
        if q.Tx != nil {
            var item T
            query := fmt.Sprintf("SELECT * FROM %s WHERE %s", 
                                    item.TableName(),
                                    constraint.constraint, 
                                )
            
            err := q.Tx.Select(&q.Rows, query, constraint.values...)
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

func SelectOffset[T QueryTypes](args map[string]interface{}, limit int, skip int, sort_by string, order string, constraint Constraint) QueryHandlerFunc[T] {
    return func(q *Query[T]) error {
        if q.Tx != nil {
            var item T
            var err error
            q.Rows = make([]T, 0, len(q.A))
            number, err := constraint.GetFinalPlaceholder(); if err != nil {
                  number = 1
            }
            if len(args) == 0 {
              query := fmt.Sprintf(
                  "SELECT * FROM %s WHERE %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
                  item.TableName(),
                  constraint.constraint,
                  sort_by,
                  order,
                  number + 1,
                  number + 2,
              )
              full_args := append(constraint.values, limit, skip)
              fmt.Println(full_args...)
              fmt.Println(query)
              total, err := GetTotalCount(q, constraint); if err != nil {
                  return err
              }
              q.Total = total
              err = q.Tx.Select(&q.Rows, query, full_args...)
            } else {
              filters , values := util.GenerateFilterString(args, number, limit, skip)
              query := fmt.Sprintf(
                  "SELECT * FROM %s WHERE %s AND %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
                  item.TableName(),
                  filters,
                  constraint.constraint,
                  sort_by,
                  order,
                  number + 2,
                  number + 3,
              )
              full_args := append(constraint.values, values...)
              //count_args := append(constraint.values, values...)
              total, err := GetTotalCountFilter(q, constraint, full_args, filters); if err != nil {
                return err
              }
              q.Total = total
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

func Update[T QueryTypes](constraint Constraint) QueryHandlerFunc[T] {
    return func(q *Query[T]) error {
        if q.A == nil {
            return fmt.Errorf("No items provided to update")
        }
        
        if q.Tx == nil {
            return fmt.Errorf("Update operation on slice should have an active transaction")
        }
        
        q.Rows = make([]T, 0, len(q.A))

        for index, item := range q.A {
            number, err := constraint.GetFinalPlaceholder(); if err != nil {
                return err
            }
            fields, err := util.FieldsAndPlaceholders[T](item, number); if err != nil {
                return err
            }
            args, err := util.GetArgs[T](item); if err != nil {
              return err
            }
            query := fmt.Sprintf("UPDATE %s SET %s WHERE %s RETURNING *", 
                                 item.TableName(),
                                 fields,
                                 constraint.constraint)
            var updatedRow T
            new_args := append(constraint.values, args...)
            err = q.Tx.QueryRowx(query, new_args...).StructScan(&updatedRow); if err != nil {
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


func GetTotalCount[T QueryTypes](q *Query[T], constraint Constraint) (int, error) {
    var item T
    var totalCount int
    countQuery := fmt.Sprintf(
        "SELECT COUNT(*) FROM %s WHERE %s",
        item.TableName(),
        constraint.constraint,
    )
    err := q.Tx.Get(&totalCount, countQuery, constraint.values...)
    if err != nil {
        return 0, fmt.Errorf("failed to get total count: %w", err)
    }
    
    return totalCount, nil
}

func GetTotalCountFilter[T QueryTypes](q *Query[T], constraint Constraint, args []interface{}, filters string) (int, error) {
    var item T
    var totalCount int
    countQuery := fmt.Sprintf(
        "SELECT COUNT(*) FROM %s WHERE %s AND %s",
        item.TableName(),
        filters,
        constraint.constraint,
    )
    countArgs := append(constraint.values, args...)
    
    err := q.Tx.Get(&totalCount, countQuery, countArgs...)
    if err != nil {
        return 0, fmt.Errorf("failed to get total count: %w", err)
    }
    
    return totalCount, nil
}
