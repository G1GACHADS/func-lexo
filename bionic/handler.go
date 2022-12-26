package bionic

import (
	"github.com/G1GACHADS/func-lexo/internal/bionicreader"
	"github.com/gofiber/fiber/v2"
)

type ConvertRequest struct {
	Content  string `json:"content"`
	Fixation int    `json:"fixation"`
	Saccade  int    `json:"saccade"`
}

func Convert(c *fiber.Ctx) error {
	// Image input and OCR Logic goes here

	var request ConvertRequest

	if err := c.BodyParser(&request); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	convertedText, err := bionicreader.Convert(request.Content, request.Fixation, request.Saccade)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "The Bionic Reading API is currently unavailable, Try again later.",
		})
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)

	return c.Status(fiber.StatusOK).SendString(convertedText)
}
