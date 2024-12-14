package queries

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
  ID    string `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

func (u User) TableName() string {
  return "users"
}

func (u User) Id() interface{} {
  return u.ID
}


func ConnectInMemory() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitializeDatabase(db *sqlx.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	);`
	_, err := db.Exec(createTableQuery)
	return err
}

func TestUserTableInitialization(t *testing.T) {
	db, err := ConnectInMemory()
	if err != nil {
		t.Fatalf("Failed to connect to in-memory database: %v", err)
	}
	defer db.Close()

	
  err = InitializeDatabase(db)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

  tx := db.MustBegin()
  defer tx.Rollback()
	
  testUser := User{Name: "John Doe", Email: "john.doe@example.com"}

  _, err = InsertQuery(tx, nil, testUser); if err != nil {
    t.Fatalf("Failed to insert test user : %v", err)
  }

  var retrievedUser User
	err = tx.Get(&retrievedUser, `SELECT * FROM users WHERE email = ?`, testUser.Email)
	if err != nil {
		t.Fatalf("Failed to retrieve test user: %v", err)
	}

	if retrievedUser.Name != testUser.Name || retrievedUser.Email != testUser.Email {
		t.Fatalf("Retrieved user data mismatch. Got: %+v, Expected: %+v", retrievedUser, testUser)
	}
}
