package main

import "hello/internal/router"

func main() {
	err := router.StartRouter()
	if err != nil {
		panic(err)
	}
}
