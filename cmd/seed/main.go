package main

import (
	"gin_stuff/internals/config"
	"gin_stuff/internals/database"
	"log"
)

func main() {
	err := perform()
	if err != nil {
		log.Printf("error performing tasks %v\n", err)
	}
}

func perform() error {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("")
	}
	dbConfig := database.DBConfig{
		MaxIdleConnections: conf.GetInt("database.max_db_conns"),
		MaxOpenConnections: conf.GetInt("database.max_open_conns"),
		MaxIdleTime:        conf.GetDuration("database.max_idle_time"),
	}
	dbInstance, err := database.OpenDB(conf.GetString("database.uri"), dbConfig)
	if err != nil {
		log.Fatalf("error open connection to database %v\n", err)
	}
	return dbInstance.Ping()
}
