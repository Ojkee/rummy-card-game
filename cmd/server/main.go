package main

import (
	"rummy-card-game/src/network_server"
)

func main() {
	server := network_server.NewServer(2)
	server.Init("8080")
}
