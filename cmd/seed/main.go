package main

import (
	"context"
	"encoding/json"
	"gin_stuff/internals/app"
	"gin_stuff/internals/config"
	"gin_stuff/internals/data"
	"io"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
)

func main() {
	config.LoadConfig()
	conn, err := app.OpenDB()
	if err != nil {
		log.Fatalf("fail to open db %s", err)
	}
	defer conn.Close(context.Background())

	query := data.New(conn)
	user, err := query.InsertUser(context.Background(), data.InsertUserParams{
		Username: "user_01",
		// Password: "Novelism_0123"
		PasswordHash: "$2a$14$B1G2aDfd50fwZ1YnyCvo.u/NCVBOUTSa13HwtT6Vti7athA8FOeHC",
		Email:        "user_01@novelism.io",
		FirstName:    pgtype.Text{String: "FIrst", Valid: true},
		LastName:     pgtype.Text{String: "Last", Valid: true},
		Status:       data.NullUserStatus{UserStatus: data.UserStatusActive, Valid: true},
		Verified:     pgtype.Bool{Bool: true, Valid: true},
		Gender:       pgtype.Text{String: "male", Valid: true},
	})
	if err != nil {
		log.Printf("Error while inserting user: %s", err.Error())
	}
	log.Printf("1 user inserted: %d: %s", user.ID, user.Username)

	file, err := os.Open("cmd/seed/sample_books.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	bytevalue, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	var books []data.BulkInsertBooksParams
	if err := json.Unmarshal(bytevalue, &books); err != nil {
		log.Fatal(err)
	}
	rowCount, err := query.BulkInsertBooks(context.Background(), books)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Printf("Inserted: %d\n", rowCount)
	}
}
