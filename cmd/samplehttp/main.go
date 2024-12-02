package main

import "github.com/eser/go-service/pkg/samplehttp"

func main() {
	err := samplehttp.Run()
	if err != nil {
		panic(err)
	}
}
