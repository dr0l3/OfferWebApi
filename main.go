package main

import "github.com/dr0l3/offerwebapi/offerrecords"

func main() {
	databaseconnection := offerrecords.NewPostGresConnnection("recordsserver")
	server := offerrecords.NewServer(databaseconnection)
	server.Run(":8080")
}
