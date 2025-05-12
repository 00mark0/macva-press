package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbSource := os.Getenv("DB_URL_TEST")

	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal("Cannot connect to db!:", err)
	}
	log.Println("Connected to db.")

	testQueries = New(conn)

	os.Exit(m.Run())
}
