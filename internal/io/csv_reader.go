package io

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/fadeldnswr/fo-performance-engine.git/internal/model"
)

// Define struct for CSV reading process
type CSVReadOptions struct {
	Delimiter rune
	DecimalComma bool
}

// Define function to read CSV file
func ReadLinksCSV(path string, opt CSVReadOptions) ([]model.LinkInput, []model.RowError, error) {
	// Open file path
	file, err := os.Open(path)

	// Check if the file is valid
	if err != nil {
		return nil, nil, errors.New("Failed to open CSV file: " + err.Error())
	}
	defer file.Close()

	// Read and parse CSV content
	reader := csv.NewReader(file)
	if opt.Delimiter != 0 {
		reader.Comma = opt.Delimiter
	}

	// Read all records
	records, err := reader.Read()
	if err != nil {
		return nil, nil, errors.New("Failed to read CSV file: " + err.Error())
	}

	// Map column headers to indices
	col := make(map[string]int, len(records))
	for i, h := range records {
		col[strings.ToLower(strings.TrimSpace(h))] = i
	}

	// Check required columns 
	var schemaErrors []model.RowError
	for _, req := range RequiredColumns {
		if _, ok := col[req]; !ok {
			schemaErrors = append(schemaErrors, model.RowError{
				Row: 0, 
				Field: req, 
				Message: "Missing required column",
			})
		}
	}
	// Check the len of the schema errors
	if len(schemaErrors) > 0 {
		return nil, schemaErrors, errors.New("CSV schema validation failed")
	}
	// Parse float value
	parseFloat := func(s string) (float64, error) {
		s = strings.TrimSpace(s)
		if s == "" {
			return 0, nil
		}
		if opt.DecimalComma {
			s = strings.ReplaceAll(s, ",", ".")
		}
		return strconv.ParseFloat(s, 64)
	}

	// Parse integer value
	parseInt := func(s string) (int, error) {
		s = strings.TrimSpace(s)
		if s == "" {
			return 0, nil
		}
		return strconv.Atoi(s)
	}

	// Define variable for output data and errors
	var output []model.LinkInput
	var rowErrs []model.RowError
	rowIndex := 0 // Row index starts from 0 (header row)
	
	// Iterate through each record
	for {
		rec, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, nil, errors.New("CSV has reached the end of the file")
		}
		rowIndex++

		// Generate function to capture row errors
		get := func(name string) string { return rec[col[name]]}

		// Create new link input to capture id and scenarios
		var link model.LinkInput
		link.LinkID = get("link_id")
		link.Scenario = get("scenario")

		// Parse and assign float and int values
		if link.TXPowerDbm, err = parseFloat(get("tx_power_dbm")); err != nil {
			rowErrs = append(rowErrs, model.RowError{Row: rowIndex, Field: "tx_power_dbm", Message: "Not a number"})
			continue
		}
		if link.RXSensitivityDbm, err = parseFloat(get("rx_sensitivity_dbm")); err != nil {
			rowErrs = append(rowErrs, model.RowError{Row: rowIndex, Field: "rx_sensitivity_dbm", Message: "Not a number"})
			continue
		}
		if link.SystemMarginDb, err = parseFloat(get("system_margin_db")); err != nil {
			rowErrs = append(rowErrs, model.RowError{Row: rowIndex, Field: "system_margin_db", Message: "Not a number"})
			continue
		}
		if link.FiberLengthKm, err = parseFloat(get("fiber_length_km")); err != nil {
			rowErrs = append(rowErrs, model.RowError{Row: rowIndex, Field: "fiber_length_km", Message: "Not a number"})
			continue
		}
		if link.FiberAttDbPerKm, err = parseFloat(get("fiber_att_db_per_km")); err != nil {
			rowErrs = append(rowErrs, model.RowError{Row: rowIndex, Field: "fiber_att_db_per_km", Message: "Not a number"})
			continue
		}
		if link.NSplice, err = parseInt(get("n_splice")); err != nil {
			rowErrs = append(rowErrs, model.RowError{Row: rowIndex, Field: "n_splice", Message: "Not an integer value"})
			continue
		}
		if link.NConnectors, err = parseInt(get("n_connector")); err != nil {
			rowErrs = append(rowErrs, model.RowError{Row: rowIndex, Field: "n_connector", Message: "Not an integer value"})
			continue
		}
		if link.SpliceLossDb, err = parseFloat(get("splice_loss_db")); err != nil {
			rowErrs = append(rowErrs, model.RowError{Row: rowIndex, Field: "splice_loss_db", Message: "Not a number"})
			continue
		}
		if link.ConnectorLossDb, err = parseFloat(get("connector_loss_db")); err != nil {
			rowErrs = append(rowErrs, model.RowError{Row: rowIndex, Field: "connector_loss_db", Message: "Not a number"})
			continue
		}
		if link.SplitterLossDb, err = parseFloat(get("splitter_loss_db")); err != nil {
			rowErrs = append(rowErrs, model.RowError{Row: rowIndex, Field: "splitter_loss_db", Message: "Not a number"})
			continue
		}
		output = append(output, link)
	}
	return output, rowErrs, nil
}