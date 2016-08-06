package offerrecords

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type DatabaseConnection struct {
	Db           *sql.DB
	DatabaseName string
}

func NewMySQLConnection(name string) *DatabaseConnection {
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASSWORD")
	databaseaddress := os.Getenv("DBADDRESS")
	databasename := os.Getenv("DBNAME")
	openstring := user + ":" + password + "@/tcp(" + databaseaddress + ")" + databasename + "?parseTime=true"

	db, err := sql.Open("mysql", openstring)
	if err != nil {
		log.Print("Error in connection to db: " + err.Error())
	}
	return &DatabaseConnection{db, name}
}

func NewPostGresConnnection(name string) *DatabaseConnection {
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASSWORD")
	databaseaddress := os.Getenv("DATBASEADDRESS")
	databasename := os.Getenv("DBNAME")

	connectionstring := fmt.Sprintf("user=%[1]s dbname=%[2]s password=%[3]s host=%[4]s sslmode=disable", user, databasename, password, databaseaddress)
	log.Print("Connectionstring: " + connectionstring)
	fmt.Println(connectionstring)
	db, err := sql.Open("postgres", connectionstring)
	if err != nil {
		log.Print("Error in connection to db: " + err.Error())
	}
	return &DatabaseConnection{db, name}
}
