package main

import (
	"fmt"
	"os"

	"github.com/dr0l3/offerwebapi/offerrecords"
)

func main() {
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASSWORD")
	databaseaddress := os.Getenv("DBADDRESS")
	databasename := os.Getenv("DBNAME")

	if user == "" || password == "" || databaseaddress == "" || databasename == "" {
		fmt.Println("Neccesary enviromentvariables not set. export DBUSER, DBPASSWORD, DBADDRESS and DBNAME")
		return
	}

	databaseconnection := offerrecords.NewPostGresConnnection("recordsserver")
	server := offerrecords.NewServer(databaseconnection)
	server.Run(":8080")
}
