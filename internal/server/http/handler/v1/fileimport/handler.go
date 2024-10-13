package fileimport

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/errorhandler"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/pkg/errors"
)

var expectedCSVHeaders = []string{"Gebiet", "Straße", "BaumNr.", "Gattung/Art", "Hochwert", "Rechtswert", "Pflanzjahr"}

const (
	fromEPSG = 31467
	toEPSG   = 4326
)

// @Summary		Import trees from a CSV file
// @Description	Import trees from a CSV file
// @Id				import-trees-from-csv
// @Tags			File Import
// @Success		204
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/import/csv [post]
// @Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
// @Param			file			formData	file	true	"CSV file to import"
func ImportTreesFromCSV(svc service.TreeService) fiber.Handler {
	start := time.Now()
	return func(c *fiber.Ctx) error {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to retrieve file",
			})
		}

		if !isCSVFile(fileHeader) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Uploaded file is not a valid CSV file",
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
				slog.Error("Error closing file:", "error", err)
			}
		}(file)

		// Validate that the file is a valid CSV
		if err := validateCSV(file); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid CSV format",
				"details": err.Error(),
			})
		}

		// Reset file cursor to the beginning before processing
		if _, err := file.Seek(0, 0); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to reset file cursor",
			})
		}

		if err := processCSVFile(c.Context(), file, svc); err != nil {
			slog.Error("Error processing csv file:", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process csv file",
				"err":   err.Error(),
			})
		}

		elapsed := time.Since(start)
		fmt.Println("Importing trees from CSV took", elapsed)
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func isCSVFile(fileHeader *multipart.FileHeader) bool {
	// Check the file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	slog.Debug("File extension", "ext", ext, "content-type", fileHeader.Header.Get("Content-Type"))
	return ext == ".csv" || fileHeader.Header.Get("Content-Type") == "text/csv"
}

func validateCSV(file multipart.File) error {
	reader := csv.NewReader(file)
	_, err := reader.ReadAll()
	return err
}

func processCSVFile(ctx context.Context, file multipart.File, svc service.TreeService) error {
	r := csv.NewReader(file)
	r.LazyQuotes = true
	header, err := r.Read()
	if err != nil {
		slog.Error("Failed to read CSV", "error", err)
		return errorhandler.HandleError(err)
	}

	// Create the header map
	headerIndexMap, err := createHeaderIndexMap(header)
	if err != nil {
		slog.Error("Error creating header index map", "error", err)
		return errorhandler.HandleError(err)
	}

	transformer, err := NewGeoTransformer(fromEPSG, toEPSG)
	if err != nil {
		return errors.Wrap(err, "error creating transformer")
	}

	var trees []*domain.TreeImport
	for i := range utils.NumberSequence(1) {
		row, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		tree, err := parseRowToTree(i, row, headerIndexMap)
		if err != nil {
			return err
		}
		trees = append(trees, tree)
	}

	geoPoints := utils.Map(trees, func(tree *domain.TreeImport) GeoPoint {
		return GeoPoint{X: tree.Latitude, Y: tree.Longitude}
	})

	transformedPoints, err := transformer.TransformBatch(geoPoints)
	if err != nil {
		return errors.Wrap(err, "failed to transform coordinates from Gauß-Krüger to WGS84")
	}

	for i, tree := range trees {
		tree.Latitude = transformedPoints[i].X
		tree.Longitude = transformedPoints[i].Y
	}

	if err := svc.ImportTree(ctx, trees); err != nil {
		return errors.Wrap(err, "failed to import trees")
	}

	return nil
}

func createHeaderIndexMap(headers []string) (map[string]int, error) {
	headerIndexMap := make(map[string]int)

	// Loop through the headers and create a map of their indices
	for i, header := range headers {
		headerIndexMap[header] = i
	}

	// Check if all expected headers are present
	for _, expectedHeader := range expectedCSVHeaders {
		if _, ok := headerIndexMap[expectedHeader]; !ok {
			return nil, errorhandler.HandleError(errors.New("missing expected header: " + expectedHeader))
		}
	}
	return headerIndexMap, nil
}

func parseRowToTree(rowIdx int, row []string, headerIndexMap map[string]int) (*domain.TreeImport, error) {
	areaIdx := headerIndexMap[expectedCSVHeaders[0]]
	area := row[areaIdx]
	if area == "" {
		return nil, errorhandler.HandleError(errors.New("invalid 'Gebiet' value at row: " + strconv.Itoa(rowIdx)))
	}

	// Read 'Species' from the row using the header index map
	streetIdx := headerIndexMap[expectedCSVHeaders[1]]
	street := row[streetIdx]

	if street == "" {
		return nil, errorhandler.HandleError(errors.New("invalid 'Straße' value at row: " + strconv.Itoa(rowIdx)))
	}

	treeNumberIdx := headerIndexMap[expectedCSVHeaders[2]]
	treeNumber := row[treeNumberIdx]

	if treeNumber == "" {
		return nil, errorhandler.HandleError(errors.New("invalid 'BaumNr.' value at row: " + strconv.Itoa(rowIdx)))
	}

	// Read 'Species' from the row using the header index map
	speciesIdx := headerIndexMap[expectedCSVHeaders[3]]
	species := row[speciesIdx]

	// if species == "" {
	// 	return nil, errorhandler.HandleError(errors.New("invalid 'Gattung/Art' value at row: " + strconv.Itoa(rowIdx)))
	// }

	latitudeIdx := headerIndexMap[expectedCSVHeaders[4]]
	longitudeIdx := headerIndexMap[expectedCSVHeaders[5]]

	latitude, err := strconv.ParseFloat(strings.ReplaceAll(row[latitudeIdx], ",", "."), 64)
	if err != nil {
		slog.Debug("Error", "error", err, "row", row, "latitudeIdx", latitudeIdx, "row[latitudeIdx]", row[latitudeIdx])
		return nil, errorhandler.HandleError(errors.New("invalid 'Hochwert' value at row: " + strconv.Itoa(rowIdx)))
	}

	longitude, err := strconv.ParseFloat(strings.ReplaceAll(row[longitudeIdx], ",", "."), 64)
	if err != nil {
		slog.Debug("Error", "error", err, "row", row, "longitudeIdx", longitudeIdx, "row[longitudeIdx]", row[longitudeIdx])
		return nil, errorhandler.HandleError(errors.New("invalid 'Rechtswert' value at row: " + strconv.Itoa(rowIdx)))
	}

	plantingYearIdx := headerIndexMap[expectedCSVHeaders[6]]
	plantingYear, err := strconv.Atoi(row[plantingYearIdx])
	if err != nil {
		return nil, errorhandler.HandleError(errors.New("invalid 'Pflanzjahr' value at row: " + strconv.Itoa(rowIdx)))
	}

	tree := &domain.TreeImport{
		Area:         area,
		Number:       treeNumber,
		Species:      species,
		Latitude:     latitude,
		Longitude:    longitude,
		PlantingYear: int32(plantingYear),
		Street:       street,
	}

	return tree, nil
}
