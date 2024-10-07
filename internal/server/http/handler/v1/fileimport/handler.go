package fileimport

import (
	"context"
	"encoding/csv"
	"errors"
	"github.com/gofiber/fiber/v2"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/errorhandler"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"log"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var expectedCSVHeaders = []string{"Area", "Street", "TreeNumber", "Species", "Latitude", "Longitude", "PlantingYear"}
var headerRowIndex = 0

func ImportTreesFromCSV(svc service.TreeService) fiber.Handler {
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
func isCSVFile(fileHeader *multipart.FileHeader) bool {
	// Check the file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".csv" {
		return false
	}

	// Check the MIME type (could be "text/csv" or "application/vnd.ms-excel")
	if fileHeader.Header.Get("Content-Type") != "text/csv" && fileHeader.Header.Get("Content-Type") != "application/vnd.ms-excel" {
		return false
	}

	return true
}
func validateCSV(file multipart.File) error {
	reader := csv.NewReader(file)

	headers, err := reader.Read()
	if err != nil {
		return errorhandler.HandleError(errors.New("failed to read CSV headers"))
	}

	if !validateHeaders(headers) {
		return errorhandler.HandleError(errors.New("CSV headers do not match the expected format"))
	}

	for {
		record, err := reader.Read()
		if err != nil {
			break // End of file reached or error (if EOF)
		}

		if len(record) != len(expectedCSVHeaders) {
			return errorhandler.HandleError(errors.New("inconsistent number of fields in CSV rows"))
		}
	}

	return nil
}

// validateHeaders compares the CSV headers to the expected headers
func validateHeaders(headers []string) bool {
	if len(headers) != len(expectedCSVHeaders) {
		return false
	}
	for i, header := range headers {
		if header != expectedCSVHeaders[i] {
			return false
		}
	}
	return true
}
func processCSVFile(ctx context.Context, file multipart.File, svc service.TreeService) error {

	r := csv.NewReader(file)
	r.LazyQuotes = true

	rows, err := r.ReadAll()
	if err != nil {
		slog.Error("Failed to read CSV: %v", err)
		return err
	}
	if len(rows) < headerRowIndex+2 {
		return errorhandler.HandleError(errors.New("CSV file has no data"))
	}

	//transformer, err := proj.NewEPSGTransformer(31467, 4326)
	//if err != nil {
	//	slog.Error("Error creating transformer: %v", err)
	//}

	// Create the header map
	headerIndexMap, err := createHeaderIndexMap(expectedCSVHeaders)
	if err != nil {
		slog.Error("Error:", err)
		return err
	}

	var trees []*domain.Tree
	for i, row := range rows[headerRowIndex+1:] {
		rowIndex := headerRowIndex + i + 1
		if len(row) < len(expectedCSVHeaders) {
			return errorhandler.HandleError(errors.New("invalid CSV format at row " + strconv.Itoa(rowIndex)))
		}

		//tree, err := parseRowToTree(rowIndex, row, headerIndexMap, &transformer)
		tree, err := parseRowToTree(rowIndex, row, headerIndexMap)
		if err != nil {
			slog.Error("Failed to parse row %d: %v", i+2, err)
			return err
		}

		trees = append(trees, tree)
	}
	err = svc.ImportTree(ctx, trees)
	if err != nil {
		return err
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

//func parseRowToTree(rowIdx int, row []string, headerIndexMap map[string]int, transformer *proj.Transformer) (*domain.Tree, error) {
func parseRowToTree(rowIdx int, row []string, headerIndexMap map[string]int) (*domain.Tree, error) {
	areaIdx := headerIndexMap[expectedCSVHeaders[0]]
	area := row[areaIdx]
	if area == "" {
		return nil, errorhandler.HandleError(errors.New("invalid 'Area' value at row: " + strconv.Itoa(rowIdx)))
	}

	streetIdx := headerIndexMap[expectedCSVHeaders[1]]
	street := row[streetIdx]
	if street == "" {
		return nil, errorhandler.HandleError(errors.New("invalid 'Street' value at row: " + strconv.Itoa(rowIdx)))
	}

	treeNumberIdx := headerIndexMap[expectedCSVHeaders[2]]
	treeNumber := row[treeNumberIdx]
	if treeNumber == "" {
		return nil, errorhandler.HandleError(errors.New("invalid 'TreeNumber' value at row: " + strconv.Itoa(rowIdx)))
	}

	// Read 'Species' from the row using the header index map
	speciesIdx := headerIndexMap[expectedCSVHeaders[3]]
	species := row[speciesIdx]
	if species == "" {
		return nil, errorhandler.HandleError(errors.New("invalid 'Species' value at row: " + strconv.Itoa(rowIdx)))
	}

	latitudeIdx := headerIndexMap[expectedCSVHeaders[4]]
	longitudeIdx := headerIndexMap[expectedCSVHeaders[5]]

	latitude, err := strconv.ParseFloat(row[latitudeIdx], 64)
	if err != nil {
		return nil, errorhandler.HandleError(errors.New("invalid 'Latitude' value at row: " + strconv.Itoa(rowIdx)))
	}

	longitude, err := strconv.ParseFloat(row[longitudeIdx], 64)
	if err != nil {
		return nil, errorhandler.HandleError(errors.New("invalid 'Longitude' value at row: " + strconv.Itoa(rowIdx)))
	}

	//TODO: Convert Gauß-Krüger coordinates to WGS84 using the transformer
	//points := []proj.Coord{
	//	proj.XY(latitude, longitude),
	//}

	//if err := transformer.Transform(points); err != nil {
	//	return nil, errorhandler.HandleError(errors.New("failed to transform coordinates from Gauß-Krüger to WGS84 at row: " + strconv.Itoa(rowIdx)))
	//}

	plantingYearIdx := headerIndexMap[expectedCSVHeaders[6]]
	plantingYear, err := strconv.Atoi(row[plantingYearIdx])
	if err != nil {
		return nil, errorhandler.HandleError(errors.New("invalid 'PlantingYear' value at row: " + strconv.Itoa(rowIdx)))
	}

	tree := &domain.Tree{
		Number:       treeNumber,
		Species:      species,
		Latitude:     latitude,  //points[0].Y, // WGS84 Latitude
		Longitude:    longitude, //points[0].X, // WGS84 Longitude
		PlantingYear: int32(plantingYear),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return tree, nil
}
