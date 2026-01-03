package calc

import (
	"fmt"
	"sort"

	"github.com/fadeldnswr/fo-performance-engine.git/internal/model"
)

// Define struct for runner options
type RunnerOptions struct {
	EnableRTB bool
	// RTB defaults
	BitrateGbps      float64
	TxRiseTimeNs     float64
	RxRiseTimeNs     float64
	DispersionPerKm  float64
}
// Define function to run calculations on link inputs
func Compute(link model.LinkInput, opt RunnerOptions) (model.LinkOutput, error) {
	// Loss precompute and breakdown
	fiberLossDb := link.FiberLengthKm * link.FiberAttDbPerKm
	connTotalDb := float64(link.NConnectors) * link.ConnectorLossDb
	spliceTotalDb := float64(link.NSplice) * link.SpliceLossDb

	// Call LPB calculation
	lpbInput := LPBInputs{
		TxPowerDbm: link.TXPowerDbm,
		RxSensitivityDbm: link.RXSensitivityDbm,
		FiberAttDbPerKm: link.FiberAttDbPerKm,
		ConnLossDb: connTotalDb,
		SpliceLossDb: spliceTotalDb,
		SystemMarginDb: link.SystemMarginDb,
		LinkLengthKm: link.FiberLengthKm,
		OtherLossDb: link.OtherLossDb,
		SplitterLossDb: link.SplitterLossDb,
	}
	lpbOutput, err := CalculateLPB(lpbInput)
	if err != nil {
		return model.LinkOutput{}, fmt.Errorf("An error has occurred: %v", err.Error())
	}

	// Prepare result structure
	res := model.LinkOutput {
		LinkID:   link.LinkID,
		Scenario: link.Scenario,
		TotalLossDb: lpbOutput.TotalLossDb,
		RxPowerDbm:  lpbOutput.RxPowerDbm,
		MarginDb:    lpbOutput.MarginDb,
		LPBStatus:   lpbOutput.Status,
		FiberLossDb: fiberLossDb,
		ConnectorTotalDb: connTotalDb,
		SpliceTotalDb: spliceTotalDb,
	}

	// Explainability placeholders
	type contribute struct {
		name string
		value float64
	}
	contributors := []contribute{
		{"fiber_loss_db", fiberLossDb},
		{"connector_total_db", connTotalDb},
		{"splice_total_db", spliceTotalDb},
		{"splitter_loss_db", link.SplitterLossDb},
		{"splice_loss_db", link.SpliceLossDb},
	}
	// Sort contributors by value descending
	sort.Slice(contributors, 
		func(i, j int) bool { 
			return contributors[i].value > contributors[j].value 
		})
	res.TopContributor1 = contributors[0].name
	res.TopContributor2 = contributors[1].name
	res.TopContributor3 = contributors[2].name

	// RTB calculation if enabled
	if opt.EnableRTB {
		rtbIn := RTBInputs{
			BitrateGbps:      opt.BitrateGbps,
			TxRiseTimeNs:     opt.TxRiseTimeNs,
			RxRiseTimeNs:     opt.RxRiseTimeNs,
			FiberLengthKm:    link.FiberLengthKm,
			DispersionPerKm:  opt.DispersionPerKm,
		}
		rtbOut, err := CalculateRTB(rtbIn)
		if err != nil {
			return model.LinkOutput{}, fmt.Errorf("An error has occurred: %v", err.Error())
		}
		// Define result fields
		res.SystemRiseTimeNs = rtbOut.TotalRiseTimeNs
		res.AllowedRiseTimeNs = rtbOut.AllowedRiseTimeNs
		res.RTBStatus = (rtbOut.Status == "PASS")
	}
	return res, nil
}