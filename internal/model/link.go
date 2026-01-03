package model

// Define link input contract data
type LinkInput struct {
	// Identifiers
	LinkID string
	Scenario string

	// Transmitter and receiver parameters
	TXPowerDbm float64
	RXSensitivityDbm float64
	SystemMarginDb float64

	// Fiber and component parameters
	FiberLengthKm float64
	FiberAttDbPerKm float64

	// Component losses and counts
	NSplice int
	SpliceLossDb float64
	NConnectors int
	ConnectorLossDb float64
	SplitterLossDb float64
	OtherLossDb float64
}

// Define link output contract data
type LinkOutput struct {
	// Identifiers
	LinkID string
	Scenario string

	// Computed loss
	FiberLossDb float64
	SpliceTotalDb float64
	ConnectorTotalDb float64
	TotalLossDb float64

	// Link power budget
	RxPowerDbm float64
	MarginDb float64
	LPBStatus string

	// Rise time budget
	SystemRiseTimeNs float64
	AllowedRiseTimeNs float64
	RTBStatus bool

	// Explainability
	TopContributor1 string
	TopContributor2 string
	TopContributor3 string
}