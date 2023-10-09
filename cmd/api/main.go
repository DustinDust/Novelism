package main

import (
	_ "github.com/lib/pq"
)

func main() {
	app := NewApplication()
	app.EchoInstance.Logger.Fatal(app.EchoInstance.Start(app.Config.GetString("general.server")))
}
