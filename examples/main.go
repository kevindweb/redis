package main

import (
	"cache/pkg/client"
	"cache/pkg/server"
	"cache/pkg/util"
	"fmt"

	"github.com/google/uuid"
)

func main() {
	defaultParameters()

	customParameters()
}

func defaultParameters() {
	fmt.Println("Starting default server and client")
	client, server, err := util.StartDefaultClientServer()
	if err != nil {
		panic(fmt.Errorf("error starting default client & server: %w", err))
	}

	defer cleanup(client, server)
	setOperation(client)
}

func customParameters() {
	host := "0.0.0.0"
	port := 8080
	network := "tcp"
	serverOptions := server.Options{
		Host:    host,
		Port:    port,
		Network: network,
	}
	customServer, err := server.New(serverOptions)
	if err != nil {
		panic(fmt.Errorf("error creating custom server: %w", err))
	}

	fmt.Printf("Starting custom server at %s\n", customServer.Address)

	go func() {
		err := customServer.Start()
		if err != nil {
			panic(fmt.Errorf("error starting custom server: %w", err))
		}
	}()

	clientOptions := client.Options{
		Host:    host,
		Port:    port,
		Network: network,
	}
	customClient, err := client.New(clientOptions)
	if err != nil {
		panic(fmt.Errorf("error creating custom client: %w", err))
	}

	customClient.Start()

	defer cleanup(customClient, customServer)
	setOperation(customClient)
}

func setOperation(client *client.Client) {
	fmt.Println("Running example set operation")
	key := uuid.New().String()
	val := uuid.New().String()
	err := client.Set(key, val)
	if err != nil {
		panic(fmt.Errorf("error setting key (%s) and val (%s): %w", key, val, err))
	}
}

func cleanup(client *client.Client, server *server.Server) {
	if err := server.Stop(); err != nil {
		panic(fmt.Errorf("error stopping server: %w", err))
	}
	if err := client.Stop(); err != nil {
		panic(fmt.Errorf("error stopping client: %w", err))
	}
}
