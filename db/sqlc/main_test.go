package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ridwan414-hub/simplebank/utils"
)


var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M){
	var err error

	config,err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config: ",err)
	}

	testDB,err = sql.Open(config.DBDriver,config.DBSource)

	if err!=nil{
		log.Fatal("Cannot connect to the db: ",err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}