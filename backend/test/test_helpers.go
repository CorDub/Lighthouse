package th

import (
	"database/sql"
	"os"
	"testing"
	"context"
	"net/http/httptest"
	"net/http"
	"encoding/json"
	"bytes"

	"Lighthouse/internal/database"
	"Lighthouse/internal/auth"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

type ErrorResponse struct{
	Error string `json:"error"`
}

func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	godotenv.Load("../.env")
	dbTestUrl := os.Getenv("DB_TEST_URL")
	db, err := sql.Open("postgres", dbTestUrl)
	if err != nil {
		t.Fatalf("Error opening the test db: %s", err)
	}

	errDialect := goose.SetDialect("postgres")
	if errDialect != nil {
		t.Fatalf("Error setting the dialect for Goose: %s", errDialect)
	}

	errMigrations := goose.Up(db, "../sql/schema") 
	if errMigrations != nil {
		t.Fatalf("Error running the migrations on test db: %s", errMigrations)
	}

	t.Cleanup(func() { db.Close() })

	return db
}


func SetupTestTransaction(t *testing.T, db *sql.DB) *database.Queries {
	t.Helper()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error setting up the transaction for the test DB: %s", err)
	}

	t.Cleanup(func() { tx.Rollback() })

	return database.New(tx)
}


func AddUser(t *testing.T, tx *database.Queries, email string, password string) {
	t.Helper()

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("Couldn't hash password: %s", err)
	}

	createUserParams := database.CreateUserParams{
		Email: email,
		HashedPassword: hash,
	}

	_, err2 := tx.CreateUser(context.Background(), createUserParams) 
	if err2 != nil {
		t.Fatalf("Couldn't create user: %s", err2)
	}
}


func CreateTestWR(
	t *testing.T,
	method string, 
	url string, 
	payload any,
	) (w *httptest.ResponseRecorder, r *http.Request) {
	
	t.Helper()

	// get the payload in bytes for the reader
	var byteBody []byte
	if payload != nil {
		body, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("Could not marshal the payload: %s", err)
		}
		byteBody = body
	}	

	// create w and r
	write := httptest.NewRecorder()
	read := httptest.NewRequest(method, url, bytes.NewReader(byteBody))
	read.Header.Set("Content-Type", "application/json")

	return write, read
}


func DecodeResponse[T any](t *testing.T, w *httptest.ResponseRecorder) T {
	t.Helper()

	var result T
	if err := json.NewDecoder(w.Result().Body).Decode(&result); err != nil {
		t.Fatalf("Could not decode response body: %s", err)
	}

	return result
}