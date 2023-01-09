package server

import (
	"os"

	"github.com/G1GACHADS/func-lexo/bionic"
	"github.com/G1GACHADS/func-lexo/bionicconfig"
	"github.com/G1GACHADS/func-lexo/internal/clients"
	"github.com/G1GACHADS/func-lexo/internal/logger"
	"github.com/gofiber/fiber/v2"
)

func New(clients clients.Clients) *fiber.App {
	srv := fiber.New(fiber.Config{
		AppName: os.Getenv("AZURE_FUNCTIONAPP_NAME"),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			statusCode := fiber.StatusInternalServerError

			e, ok := err.(*fiber.Error)
			if ok {
				statusCode = e.Code
			} else {
				logger.M.Error(err)
			}

			return c.Status(statusCode).JSON(fiber.Map{
				"message": e.Message,
			})
		},
	})

	handlerBionicConverter := bionic.NewHandler(clients.AzureComputerVisionClient)
	srv.Post("/api/bionic", handlerBionicConverter.Handler)

	handlerBionicConverterRaw := bionicconfig.NewHandler()
	srv.Post("/api/bionicconfig", handlerBionicConverterRaw.Handler)

	return srv
}
