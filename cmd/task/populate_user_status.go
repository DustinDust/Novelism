package main

import (
	"gin_stuff/internals/config"
	"gin_stuff/internals/database"
	"log"

	_ "github.com/lib/pq"
)

/*
- Update user without a status to have correct status.

- This is used when `status` is added to `users“ table.
*/
func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config %v", err)
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

	statement := `
		UPDATE users SET status = 'active' WHERE status IS NULL
	`
	tx := dbInstance.MustBegin()
	res := tx.MustExec(statement)
	if row_aff, err := res.RowsAffected(); row_aff == 0 || err != nil {
		if err == nil {
			log.Println("Updated 0 row")
		}
		log.Fatalf("Something wrong while updating %v\n", err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Something wrong while updating %v\n", err)
	}
}
