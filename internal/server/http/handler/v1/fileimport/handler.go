package fileimport

import (
	"context"
	"encoding/csv"
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"log"
	"mime/multipart"
	"os"
)

func ImportTreesFromCSV(svc service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to retrieve file",
			})
		}

		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to open uploaded file",
			})
		}
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)

		if err := processCSVFile(c.Context(), file, svc); err != nil {
			log.Printf("Error processing csv file: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process csv file",
				"err":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "File uploaded and processed successfully",
		})
	}
}

func processCSVFile(ctx context.Context, file multipart.File, svc service.TreeService) error {
	tempFile, err := os.CreateTemp("", "upload-*.csv")
	if err != nil {
		return err
	}
	defer tempFile.Close()
	if _, err := tempFile.ReadFrom(file); err != nil {
		return err
	}
	tempFile.Seek(0, os.SEEK_SET)
	f, err := os.Open(tempFile.Name())
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.LazyQuotes = true

	rows, err := r.ReadAll()
	if err != nil {
		log.Printf("Failed to read CSV: %v", err)
		return err
	}

	err = svc.ImportTree(ctx, rows)
	if err != nil {
		return err
	}

	return nil
}
