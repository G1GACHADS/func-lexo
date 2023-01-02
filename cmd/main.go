package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/G1GACHADS/func-lexo/internal/clients"
	"github.com/G1GACHADS/func-lexo/internal/logger"
	"github.com/G1GACHADS/func-lexo/internal/server"
)

func main() {
	logger.Init(os.Getenv("ENVIRONMENT") == "dev")

	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	clients, err := clients.New(clients.ClientsConfig{
		AzureOCPSubscriptionKey: os.Getenv("OCP_APIM_SUBSCRIPTION_KEY"),
		AzureOCPEndpoint:        os.Getenv("OCP_APIM_ENDPOINT"),
	})
	if err != nil {
		logger.M.Fatal(err)
	}

	srv := server.New(clients)

	go func() {
		if err := srv.Listen(listenAddr); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	srv.Shutdown()
}
