package main

import (
	mysqldb "gogeek/connection"
	"gogeek/routes"

	"log"
)

func main() {

	db := mysqldb.SetupDB()
	r := routes.SetupRoutes(db)
	log.Println("Server started on: http://localhost:9000")
	r.Run(":9000")

}
