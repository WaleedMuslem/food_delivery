package main

import (
	"food_delivery/config"
	"food_delivery/server"
)

func main() {

	cfg := config.NewConfig()

	server.StartServer(cfg)

}
