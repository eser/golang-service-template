package main

import "github.com/eser/go-service/pkg/identitysvc"

func main() {
	err := identitysvc.Run()
	if err != nil {
		panic(err)
	}
}
