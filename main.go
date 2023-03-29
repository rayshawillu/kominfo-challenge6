package main

import (
	"book-challenge/routers"
)

func main() {
	var PORT = ":8080"
	routers.StartServer().Run(PORT)
}
