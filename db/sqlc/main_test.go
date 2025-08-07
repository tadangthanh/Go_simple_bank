package sqlc

import (
	"context"
	"log"
	"os"
	"testing"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var conn *pgxpool.Pool

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

// diem vao chinh cho tat ca cac bai test golang cu the trong 1 package, o day la sqlc
func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error
	conn, err = pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close()
	testQueries = New(conn)
	os.Exit(m.Run())
}
