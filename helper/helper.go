//https://docs.microsoft.com/en-us/azure/azure-sql/database/connect-query-go
package helper

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"gopkg.in/yaml.v2"
)

var db *sql.DB

type DataBaseConfig struct {
	server   string `yaml:"server"`
	user     string `yaml:"user"`
	password string `yaml:"password"`
	port     int    `yaml:"port"`
	database string `yaml:"database"`
}

func DelSplit(tosplit string, sep rune) []string {
	var fields []string

	last := 0
	for i, c := range tosplit {
		if c == '|' {
			// Found the separator, append a slice
			fields = append(fields, string(tosplit[last:i]))
			last = i + 1
		}
	}

	// Don't forget the last field
	fields = append(fields, string(tosplit[last:]))

	return fields
}

func ConnectDB() (*sql.DB, context.Context) {
	var dberr error

	yamlBytes, err := ioutil.ReadFile("./configs/database.yml")
	if err != nil {
		log.Fatal((err))
	}
	// parse the YAML stored in the byte slice into the struct
	databaseConfig := &DataBaseConfig{}
	err = yaml.Unmarshal(yamlBytes, databaseConfig)
	if err != nil {
		log.Fatal((err))
	}

	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		databaseConfig.server, databaseConfig.user, databaseConfig.password, databaseConfig.port, databaseConfig.database)

	// Create connection pool
	db, dberr = sql.Open("sqlserver", connString)
	if dberr != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	dberr = db.PingContext(ctx)
	if dberr != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")
	return db, ctx
}

func PrettyPrintStruct(i interface{}) (string, error) {
	s, err := json.MarshalIndent(i, "", "\t")
	return string(s), err
}

// politicianInfoJson, err := helper.PrettyPrintStruct(innerDataInfo)
// if err != nil {
// 	fmt.Println(politicianInfoJson)
// } else {
// 	log.WithFields(
// 		log.Fields{
// 			"parser": "Representative",
// 			"step":   "jsonConversion",
// 		},
// 	).Warn("Error converting struct to json")
// }
