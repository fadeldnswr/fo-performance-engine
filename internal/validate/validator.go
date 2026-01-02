package validate

import "github.com/fadeldnswr/fo-performance-engine.git/internal/model"

// Define struct for options used in validation
type ValidationOptions struct {
	MaxFiberAttPerDbKm float64
}

// Define function to validate link communication parameters
func ValidateLink(links []model.LinkInput, opt ValidationOptions) []model.RowError {
	// Check if max fiber attenuation is set
	if opt.MaxFiberAttPerDbKm == 0 {
		opt.MaxFiberAttPerDbKm = 1.0 // Fiber attenuation default value in dB/km
	}

	// Define slice to hold errors
	var errs []model.RowError

	// Iterate over each link to validate
	for i, link := range links {
		row := i + 1 // Row number for error reporting

		// Validate fiber attenuation
		if link.LinkID == "" {
			errs = append(errs, model.RowError{Row: row, Field: "link_id", Message: "Required"})
		}
		if link.FiberLengthKm < 0 {
			errs = append(errs, model.RowError{Row: row, Field: "fiber_length_km", Message: "Input value must be greater than zero"})
		}
		if link.FiberAttDbPerKm <= 0 || link.FiberAttDbPerKm > opt.MaxFiberAttPerDbKm {
			errs = append(errs, model.RowError{Row: row, Field: "fiber_att_db_per_km", Message: "Input value must be greated than zero"})
		}
		if link.NSplice < 0 {
			errs = append(errs, model.RowError{Row: row, Field: "n_splice", Message: "Input value must be zero or greater"})
		}
		if link.NConnectors < 0 {
			errs = append(errs, model.RowError{Row: row, Field: "n_connector", Message: "Number of connector has to be greater than zero"})
		}
		if link.ConnectorLossDb < 0 {
			errs = append(errs, model.RowError{Row: row, Field: "connector_loss_db", Message: "Connector loss has to be zero or greater"})
		}
		if link.SplitterLossDb < 0 {
			errs = append(errs, model.RowError{Row: row, Field: "splitter_loss_db", Message: "Splitter loss has to be zero or greater"})
		}
		if link.SpliceLossDb < 0 {
			errs = append(errs, model.RowError{Row: row, Field: "splice_loss_db", Message: "Splice loss has to be zero or greater"})
		}
	}
	return errs
}