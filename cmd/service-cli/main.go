package main

import "github.com/eser/go-service/pkg/app"

func main() {
	instance := app.New()

	instance.Run()
}
