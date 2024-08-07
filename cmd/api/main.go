package main

import (
	"gin_stuff/internals/app"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	app := app.NewApplication()
	if err := app.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
