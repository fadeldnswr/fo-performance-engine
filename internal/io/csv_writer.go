package io

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"

	"github.com/fadeldnswr/fo-performance-engine.git/internal/model"
)

// Define function to write results into CSV format
func WriteCSV(path string, results []model.LinkOutput, delimiter rune) error {
	// Create or overwrite the CSV file
	file, err := os.Create(path)
	if err != nil { // Check if the path is valid
		return errors.New("Failed to create CSV file: " + err.Error())
	}
	defer file.Close()

	// Write CSV headers
	write := csv.NewWriter(file)
	if delimiter != 0 {
		write.Comma = delimiter
	}

	// Define header row
	headers := []string{
		"link_id","scenario",
		"fiber_loss_db","splice_total_db","connector_total_db","total_loss_db",
		"rx_power_dbm","margin_db","lpb_status",
		"system_rise_time_ns","allowed_rise_time_ns","rtb_pass",
		"top_contributor_1","top_contributor_2","top_contributor_3",
	}
	if err := write.Write(headers); err != nil {
		return errors.New("An error has occurred while writing CSV headers: " + err.Error())
	}

	// Format and write each result row
	formatFloat := func(x float64) string { return strconv.FormatFloat(x, 'f', 6, 64) }
	for _, res := range results { // Iterate over results
		row := []string {
			res.LinkID, res.Scenario,
			formatFloat(res.FiberLossDb), formatFloat(res.SpliceTotalDb),
			formatFloat(res.ConnectorTotalDb), formatFloat(res.TotalLossDb),
			formatFloat(res.RxPowerDbm), formatFloat(res.MarginDb), res.LPBStatus,
			formatFloat(res.SystemRiseTimeNs), formatFloat(res.AllowedRiseTimeNs), 
			strconv.FormatBool(res.RTBStatus), res.TopContributor1, res.TopContributor2,
			res.TopContributor3,
		}
		if err := write.Write(row); err != nil { // Check for write errors
			return errors.New("An error has occurred while writing CSV row: " + err.Error())
		}
	}
	write.Flush()
	return write.Error()
} 