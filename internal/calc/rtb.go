package calc

import (
	"errors"
)

// Define struct for RTB inputs
type RTBInputs struct {
	BitrateGbps float64
	TxRiseTimeNs float64
	RxRiseTimeNs float64
	FiberLengthKm float64
	DispersionPerKm float64
}

// Define struct for RTB outputs
type RTBResults struct {
	TotalRiseTimeNs float64
	AllowedRiseTimeNs float64
	Status string
}

// Define function to calculate t_chromatic
func CalculateTchrom(){}

// Define function to calculate Rise Time Budget
func CalculateRTB(input RTBInputs) (RTBResults, error){
	// Check if inputs are valid
	if input.BitrateGbps <= 0 {
		return RTBResults{}, errors.New("Bitrate value does not have valid input")
	}
	// RTB Calculation logic
	dispersion := input.FiberLengthKm * input.DispersionPerKm
	systemRt := input.TxRiseTimeNs + input.RxRiseTimeNs + dispersion

	// Allowed rise time calculation: Trx = 0.7 / Bitrate
	allowedRt := 0.7 / input.BitrateGbps * 1000 // Convert to ns

	// Determine status
	status := "FAIL"
	if systemRt <= allowedRt {
		status = "PASS"
	}
	
	// Return results
	return RTBResults{
		TotalRiseTimeNs: systemRt,
		AllowedRiseTimeNs: allowedRt,
		Status: status,
	}, nil
}