package main

import (
	"gin_stuff/internals/config"
	"gin_stuff/internals/database"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

/*
- Update user without a status to have correct status.

- This is used when `status` is added to `users“ table.
*/
func main() {
	err := perform()
	if err != nil {
		log.Printf("Error performing tast %v\n", err)
	}
}

func perform() error {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config %v", err)
	}
	dbConfig := database.DBConfig{
		MaxIdleConnections: viper.GetInt("database.max_db_conns"),
		MaxOpenConnections: viper.GetInt("database.max_open_conns"),
		MaxIdleTime:        viper.GetDuration("database.max_idle_time"),
	}

	dbInstance, err := database.OpenDB(viper.GetString("database.uri"), dbConfig)
	if err != nil {
		log.Fatalf("error open connection to database %v\n", err)
	}

	statement := `
		UPDATE users SET status = 'active' WHERE status IS NULL
	`
	tx := dbInstance.MustBegin()
	defer tx.Rollback()
	res := tx.MustExec(statement)
	if row_aff, err := res.RowsAffected(); row_aff == 0 || err != nil {
		if err == nil {
			log.Println("Updated 0 row")
		} else {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
