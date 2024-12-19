package main

import "github.com/eser/go-service/pkg/eserlivesvc"

func main() {
	err := eserlivesvc.Run()
	if err != nil {
		panic(err)
	}
}
