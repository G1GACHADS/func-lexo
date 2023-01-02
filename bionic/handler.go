package bionic

import (
	"github.com/G1GACHADS/func-lexo/internal/api"
	"github.com/G1GACHADS/func-lexo/internal/bionicreader"
	"github.com/G1GACHADS/func-lexo/internal/clients"
	"github.com/G1GACHADS/func-lexo/internal/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ConvertHandler struct {
	acvClient *clients.AzureComputerVisionClient
	log       *zap.SugaredLogger
}

func NewHandler(acvClient *clients.AzureComputerVisionClient) ConvertHandler {
	logger.M.Debug("Creating Bionic Handler")
	return ConvertHandler{
		acvClient: acvClient,
		log:       logger.M.With("handler", "bionic.Handler"),
	}
}

type ConvertParams struct {
	Fixation int `form:"fixation"`
	Saccade  int `form:"saccade"`
}

type ConvertOutput struct {
	Html        string `json:"html"`
	BoundingBox string `json:"bounding_box"`
}

var supportedImageTypes = []string{
	"image/jpeg",
	"image/png",
}

func (h ConvertHandler) Handler(c *fiber.Ctx) error {
	h.log.Debug("Handling request...")
	// Image input and OCR Logic goes here
	c.Accepts(fiber.MIMEMultipartForm)

	var request ConvertParams

	h.log.Debug("Parsing body...")
	if err := c.BodyParser(&request); err != nil {
		h.log.Debug("Failed parsing body reason:\n", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if len(form.File["image"]) == 0 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Image is required")
	}

	image := form.File["image"][0]

	h.log.Debug("Opening multipart image file...")
	file, err := image.Open()
	if err != nil {
		h.log.Debug("Failed opening multipart image file reason:\n", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if !api.MimeContains(file, supportedImageTypes) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":         "Invalid image type",
			"file":            image.Filename,
			"supported_types": supportedImageTypes,
		})
	}

	h.log.Debug("Sending image into Azure Cognitive Service for OCR...")
	ocrResult, err := h.acvClient.RecognizePrintedOCR(c.Context(), file)
	if err != nil {
		h.log.Debug("Failed sending image into Azure Cognitive Service for OCR reason:\n", err.Error())
		return fiber.NewError(fiber.StatusServiceUnavailable, "The Bionic Reading API is currently unavailable, Try again later.")
	}

	mergedOCRResult := h.acvClient.MergeOCRResultLines(ocrResult)

	h.log.Debug("Converting OCR result into bionic text...")

	var output ConvertOutput

	output.Html, err = bionicreader.Convert(mergedOCRResult.MergedContent, request.Fixation, request.Saccade)
	if err != nil {
		h.log.Debug("Failed converting OCR result into bionic text reason:\n", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Failed converting OCR result into bionic text",
			"error":   err.Error(),
		})
	}

	h.log.Debug("Result:")
	h.log.Debug(output)

	return c.Status(fiber.StatusOK).JSON(output)
}
