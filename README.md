# Go Generic SQLx Queries

This project demonstrates how to use generic functions with `sqlx` to perform database operations with reduced code repetition by leveraging Go generics.

## Features

- **Generic SQL operation Functions**: Dynamically constructs an  SQL statement from any struct using reflection and `sqlx`.
  
## Setup

1. Clone the repository:

   ```bash
   go get https://github.com/nam2184/go-generic-sqlx.git
   ```

2. Add your business logic structures that correspond to the tables within model package and specify TableName() to match QueryTypes interface
3. Initialise database 

  ```go
  psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
      Host, Port, User, Password, Name)
  
  db, err := sqlx.Connect("postgres", psqlInfo)
  ```

4. You can create a query such as insert with this code for a slice of the table structure

  ```go
    // Define the SQL insert query
    tx := db.MustBegin()
    defer func() {
        if err := tx.Rollback(); err != nil {
            t.Fatal(err)
        }
    }()

    InsertQuery[model.Task, int64](tx, nil, tasks)
  ```
