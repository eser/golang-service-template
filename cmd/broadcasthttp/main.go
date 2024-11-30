package main

import "github.com/eser/go-service/pkg/broadcasthttp"

func main() {
	err := broadcasthttp.Run()
	if err != nil {
		panic(err)
	}
}
