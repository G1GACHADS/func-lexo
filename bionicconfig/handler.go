package bionicconfig

import (
	"github.com/G1GACHADS/func-lexo/internal/bionicreader"
	"github.com/G1GACHADS/func-lexo/internal/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ConvertRawHandler struct {
	log *zap.SugaredLogger
}

func NewHandler() ConvertRawHandler {
	logger.M.Debug("Creating Bionic Handler")
	return ConvertRawHandler{
		log: logger.M.With("handler", "bionic.Handler"),
	}
}

type ConvertRawParams struct {
	Text     string `json:"text" form:"text"`
	Fixation int    `json:"fixation" form:"fixation"`
	Saccade  int    `json:"saccade" form:"saccade"`
}

type ConvertRawOutput struct {
	Result      string `json:"result"`
	ResultRaw   string `json:"result_raw"`
	BoundingBox string `json:"bounding_box"`
}

func (h ConvertRawHandler) Handler(c *fiber.Ctx) error {
	var request ConvertRawParams

	if err := c.BodyParser(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := bionicreader.Convert(request.Text, request.Fixation, request.Saccade)
	if err != nil {
		h.log.Warn("Failed converting OCR result into bionic text reason:\n", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Failed converting OCR result into bionic text",
			"error":   err.Error(),
		})
	}

	// caching logic maybe

	return c.Status(fiber.StatusOK).JSON(ConvertRawOutput{
		Result:    result.Markdown,
		ResultRaw: result.Text,
	})
}
