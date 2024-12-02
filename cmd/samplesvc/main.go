package main

import "github.com/eser/go-service/pkg/samplesvc"

func main() {
	err := samplesvc.Run()
	if err != nil {
		panic(err)
	}
}
