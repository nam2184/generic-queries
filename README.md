# Go Generic SQLx Queries

This project demonstrates how to use generic functions with `sqlx` to perform PostGres database operations with reduced code repetition by leveraging Go generics, mostly for me but idk use it if u want to.

## Features

- **Generic SQL operation Functions**: Dynamically constructs an  SQL statement for PostGres from any struct using reflection and `sqlx`.
  
## Setup

1. Clone the repository:

   ```bash
   go get github.com/nam2184/generic-queries
   ```

2. Add your business logic that specify TableName() and Id() to match QueryTypes interface with db tags for sqlx
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

    queries.InsertQuery[model.Task](tx, nil, tasks)
  ```
