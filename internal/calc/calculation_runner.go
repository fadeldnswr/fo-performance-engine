package calc

import (
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
func Compute(link model.LinkInput, opt RunnerOptions) (model.LinkOutput) {
	// Call LPB calculation
	lpbInput := LPBInputs{
		TxPowerDbm: link.TXPowerDbm,
		RxSensitivityDbm: link.RXSensitivityDbm,
		FiberAttDbPerKm: link.FiberAttDbPerKm,
		ConnLossDb: float64(link.NConnectors) * link.ConnectorLossDb,
		SpliceLossDb: float64(link.NSplice) * link.SpliceLossDb,
		SystemMarginDb: link.SystemMarginDb,
		LinkLengthKm: link.FiberLengthKm,
	}
	lpbOutput, err := CalculateLPB(lpbInput)
	if err != nil {
		return model.LinkOutput{}
	}

	// Prepare result structure
	res := model.LinkOutput {
		LinkID:   link.LinkID,
		Scenario: link.Scenario,
		TotalLossDb: lpbOutput.TotalLossDb,
		RxPowerDbm:  lpbOutput.RxPowerDbm,
		MarginDb:    lpbOutput.MarginDb,
		LPBStatus:   lpbOutput.Status,
	}

	// Loss breakdown
	fiberLossDb := link.FiberLengthKm * link.FiberAttDbPerKm
	connTotalDb := float64(link.NConnectors) * link.ConnectorLossDb
	spliceTotalDb := float64(link.NSplice) * link.SpliceLossDb

	// Assign breakdown to result
	res.FiberLossDb = fiberLossDb
	res.ConnectorTotalDb = connTotalDb
	res.SpliceTotalDb = spliceTotalDb

	// Total loss breakdown match with link power budget total loss
	res.TotalLossDb = fiberLossDb + connTotalDb + spliceTotalDb + link.SplitterLossDb + link.SpliceLossDb

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
	sort.Slice(contributors, func(i, j int) bool{ return contributors[i].value > contributors[j].value })
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
			return model.LinkOutput{}
		}
		// Define result fields
		res.SystemRiseTimeNs = rtbOut.TotalRiseTimeNs
		res.AllowedRiseTimeNs = rtbOut.AllowedRiseTimeNs
		res.RTBStatus = (rtbOut.Status == "PASS")
	}
	return res
}