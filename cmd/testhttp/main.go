package main

import "github.com/eser/go-service/pkg/testhttp"

func main() {
	err := testhttp.Run()
	if err != nil {
		panic(err)
	}
}
