package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/ridwan414-hub/simplebank/api"
	db "github.com/ridwan414-hub/simplebank/db/sqlc"
)

const ( 
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@127.0.0.1:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn,err := sql.Open(dbDriver,dbSource)

	if err!=nil{
		log.Fatal("Cannot connect to the db: ",err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot start server: ",err)
	}

}