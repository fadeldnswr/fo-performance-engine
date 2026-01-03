package sweep

import (
	"errors"
	"fmt"

	"github.com/fadeldnswr/fo-performance-engine.git/internal/calc"
	"github.com/fadeldnswr/fo-performance-engine.git/internal/model"
)

// Define struct for sweep options
type SweepOptions struct {
	Runner calc.RunnerOptions
}

// Define function to apply the variations to a given input
func ApplyVariations(link model.LinkInput, v Variation, value float64) (model.LinkInput, error) {
	// Define output link as a copy of input
	out := link
	switch v.Field {
	case "engineering_margin_db", "system_margin_db":
		out.SystemMarginDb = value
	case "fiber_length_km":
		out.FiberLengthKm = value
	case "fiber_att_db_per_km":
		out.FiberAttDbPerKm = value
	case "splitter_loss_db":
		out.SplitterLossDb = value
	default:
		return link, errors.New("Unknown variation field: " + v.Field)
	}
	return out, nil
}

// Define function to perform the sweep over variations
func RunSweep(base []model.LinkInput, vars []Variation, opt SweepOptions) ([]model.LinkOutput){
	// Define slice to hold results
	results := []model.LinkOutput{}

	// Recursive generator
	var rec func(index int, current []float64)
	rec = func(index int, current []float64) {
		// Check if index is equal to length of vars
		if index == len(vars){
			scName := "base"
			for i, v := range vars {
				scName += "_" + v.Field + "=" + fmt.Sprintf("%.2f", current[i])
			}
			for _, li := range base {
				mod := li
				mod.Scenario = scName
				var err error
				for i, v := range vars {
					mod, err = ApplyVariations(mod, v, current[i])
					if err != nil { 
						panic(err)
					}
				}
				finalRes, err := calc.Compute(mod, opt.Runner)
				if err != nil { continue }
				results = append(results, finalRes)
			}
			return
		}
		for _, val := range vars[index].Values {
			next := append(current, val)
			rec(index + 1, next)
		}
	}
	rec(0, []float64{})
	return results
}