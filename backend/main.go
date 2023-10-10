package main

import (
	"nextjs/backend/api"
)

func main() {
	server := api.NewServer(".")
	server.Start(8000)
}
