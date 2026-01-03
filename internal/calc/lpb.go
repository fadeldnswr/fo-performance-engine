package calc

import (
	"errors"
)

// Define struct for LPB inputs
type LPBInputs struct {
	TxPowerDbm float64
	RxSensitivityDbm float64
	FiberAttDbPerKm float64
	ConnLossDb float64
	SpliceLossDb float64
	SplitterLossDb float64
	SystemMarginDb float64
	LinkLengthKm float64
	OtherLossDb float64
}

// Define struct for LPB outputs
type LPBResults struct {
	TotalLossDb float64
	RxPowerDbm float64
	MarginDb float64
	Status string
}

// Define function to calculate link power budget
func CalculateLPB(input LPBInputs) (LPBResults, error) {
	// Check if inputs are valid
	if input.FiberAttDbPerKm < 0 {
		return LPBResults{}, errors.New("Fiber loss does not have valid value")
	}

	// LPB Calculation logic
	fiberAttenuation := input.FiberAttDbPerKm * input.LinkLengthKm
	totalLoss := fiberAttenuation + input.ConnLossDb + input.SpliceLossDb + input.SplitterLossDb + input.OtherLossDb

	// Received power and margin calculation: Pr = Pt - Ps or Total Loss
	rxPower := input.TxPowerDbm - totalLoss
	margin := rxPower - input.RxSensitivityDbm - input.SystemMarginDb

	// Determine status
	status := "FAIL"
	if margin >= 0 {
		status = "PASS"
	}
	
	// Return results
	return LPBResults{
		TotalLossDb: totalLoss,
		RxPowerDbm: rxPower,
		MarginDb: margin,
		Status: status,
	}, nil
}