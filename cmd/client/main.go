package main

import (
	"rummy-card-game/src/network_client"
)

func main() {
	currentClient := network_client.NewClient()
	currentClient.Connect()
}
