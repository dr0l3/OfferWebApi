package offerrecords

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseConnection struct {
	Db           *sql.DB
	DatabaseName string
}

func NewConnection(name string) *DatabaseConnection {
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASSWORD")
	databasename := os.Getenv("DBNAME")
	openstring := user + ":" + password + "@/" + databasename + "?parseTime=true"

	db, _ := sql.Open("mysql", openstring)
	return &DatabaseConnection{db, name}
}
