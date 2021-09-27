package main

import (
	"log"
)

func main() {
	var (
		app       App
		err       error
	)

	if app, err = NewApplication(false); err != nil {
		log.Fatal(err)
	}

	app.Prioritise1000Payments()
}
