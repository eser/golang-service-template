package main

import "github.com/eser/go-service/pkg/broadcastsvc"

func main() {
	err := broadcastsvc.Run()
	if err != nil {
		panic(err)
	}
}
