package db_test

import (
	"database/sql"
	"fmt"
	"log"
	db "nextjs/backend/db/sqlc"
	"nextjs/backend/utils"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQuery *db.Queries

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("Could not laod env config")
	}
	conn, err := sql.Open(config.DBdriver, config.DBsource)
	if err != nil {
		log.Fatal("Could not connect to database", err)
	}

	if err := conn.Ping(); err != nil {
		fmt.Println("Couldn't ping db: ", err)
	}

	testQuery = db.New(conn)

	os.Exit(m.Run())
}
