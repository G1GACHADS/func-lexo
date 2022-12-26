package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/G1GACHADS/func-lexo/bionic"
	"github.com/gofiber/fiber/v2"
)

func main() {
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	srv := fiber.New()

	srv.Post("/api/bionic", bionic.Convert)

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
