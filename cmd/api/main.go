package main

import (
	"gin_stuff/internals/app"
	"log"
)

func main() {
	app := app.NewApplication()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
