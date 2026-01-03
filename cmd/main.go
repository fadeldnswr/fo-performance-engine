/*Fiber Optics Performance Engine CLI Application*/

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fadeldnswr/fo-performance-engine.git/internal/calc"
	foio "github.com/fadeldnswr/fo-performance-engine.git/internal/io"
	"github.com/fadeldnswr/fo-performance-engine.git/internal/model"
	"github.com/fadeldnswr/fo-performance-engine.git/internal/sweep"
	"github.com/fadeldnswr/fo-performance-engine.git/internal/validate"
)

// Main function for CLI app
func main(){
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	// Determine command args
	cmd := os.Args[1]
	switch cmd {
	case "validate":
		cmdValidate(os.Args[2:])
	case "run":
		cmdRun(os.Args[2:])
	case "sweep":
		cmdSweep(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
}

// Define function for usage of the CLI apps
func usage(){
	fmt.Println("FTTH / Fiber Optic Performance Engine")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  fo validate --in links.csv")
	fmt.Println("  fo run      --in links.csv --out results.csv [--rtb]")
	fmt.Println("  fo sweep    --in links.csv --out results.csv --vary engineering_margin_db=3,6")
}

// Define function to handle validate command
func cmdValidate(args []string){
	flagVal := flag.NewFlagSet("validate", flag.ExitOnError)
	input := flagVal.String("in", "", "input CSV")
	_ = flagVal.Parse(args)

	// Check if input is provided
	if *input == "" {
		fmt.Println("missing --in")
		os.Exit(1)
	}

	// Define slice to hold links
	links, rowErrs, err := foio.ReadLinksCSV(*input, foio.CSVReadOptions{})
	if err != nil {
		fmt.Println("An error has occurred: ", err)
		os.Exit(1)
	}

	// Check if row errors exist
	if len(rowErrs) > 0 {
		for _, e := range rowErrs {
			fmt.Println(e.Error())
		}
		os.Exit(1)
	}
	valErrs := validate.ValidateLink(links, validate.ValidationOptions{})
	if len(valErrs) > 0 {
		for _, e := range valErrs {
			fmt.Println(e.Error())
		}
		os.Exit(1)
	}
	fmt.Printf("OK — %d rows validated\n", len(links))
}

// Define function to run command
func cmdRun(args []string){
	flagRun := flag.NewFlagSet("run", flag.ExitOnError)
	input := flagRun.String("in", "", "input CSV")
	output := flagRun.String("out", "results.csv", "output CSV")

	// Rise Time Budget options
	enableRTB := flagRun.Bool("rtb", false, "enable RTB")
	bitrate := flagRun.Float64("bitrate-gbps", 2.5, "bitrate (Gbps)")
	txrt := flagRun.Float64("tx-rt-ns", 0.2, "Tx rise time (ns)")
	rxrt := flagRun.Float64("rx-rt-ns", 0.2, "Rx rise time (ns)")
	dispersion := flagRun.Float64("disp-ns-km", 0.0, "dispersion (ns/km)")

	// Parse flags
	_ = flagRun.Parse(args)

	// Check if input is provided
	if *input == "" {
		fmt.Println("missing --in")
		os.Exit(2)
	}

	// Read input CSV
	links, rowErrs, err := foio.ReadLinksCSV(*input, foio.CSVReadOptions{})

	// Check for read errors
	if err != nil {
		fmt.Println("An error has occurred: ", err.Error())
		os.Exit(1)
	}

	// Check if row errors exist
	if len(rowErrs) > 0 {
		for _, e := range rowErrs {
			fmt.Println(e.Error())
		}
		os.Exit(1)
	}

	// Process each link
	valErrs := validate.ValidateLink(links, validate.ValidationOptions{})
	if len(valErrs) > 0 {
		for _, e := range valErrs {
			fmt.Println(e.Error())
		}
		os.Exit(1)
	}

	// Define options for calculations
	opt := calc.RunnerOptions{
		EnableRTB:       *enableRTB,
		BitrateGbps:     *bitrate,
		TxRiseTimeNs:    *txrt,
		RxRiseTimeNs:    *rxrt,
		DispersionPerKm: *dispersion,
	}
	results := make([]model.LinkOutput, 0, len(links))
	for _, link := range links {
		res, err := calc.Compute(link, opt)
		if err != nil {
			fmt.Println("An error has occurred: ", err.Error())
			os.Exit(1)
		}
		results = append(results, res)
	}
	if err := foio.WriteCSV(*output, results, ','); err != nil {
		fmt.Println("An error has occurred: ", err.Error())
		os.Exit(1)
	}
	fmt.Printf("DONE — %d links written to %s\n", len(results), *output)
}

// Define function to run sweep command
func cmdSweep(args []string){
	flagSweep := flag.NewFlagSet("sweep", flag.ExitOnError)
	input := flagSweep.String("in", "", "input CSV")
	output := flagSweep.String("out", "result_sweep.csv", "output CSV")
	vary := flagSweep.String("vary", "", "variation spec (e.g. system_margin_db=3,6)")

	// Define options to enable RTB
	enableRTB := flagSweep.Bool("rtb", false, "enable RTB")
	bitrate := flagSweep.Float64("bitrate-gbps", 2.5, "bitrate (Gbps)")
	txrt := flagSweep.Float64("tx-rt-ns", 0.2, "Tx rise time (ns)")
	rxrt := flagSweep.Float64("rx-rt-ns", 0.2, "Rx rise time (ns)")
	dispersion := flagSweep.Float64("disp-ns-km", 0.0, "dispersion (ns/km)")

	// Parse flags
	_ = flagSweep.Parse(args)

	// Check if the input args are provided
	if *input == "" || *vary == "" {
		fmt.Println("missing --in or --vary")
		os.Exit(2)
	}

	// Define options for runner
	links, rowErrs, err := foio.ReadLinksCSV(*input, foio.CSVReadOptions{})
	if err != nil {
		fmt.Println("An error has occurred: ", err.Error())
		os.Exit(1)
	}

	// Check if row errors exist
	if len(rowErrs) > 0 {
		for _, e := range rowErrs {
			fmt.Println(e.Error())
		}
		os.Exit(1)
	}

	// Perform sweep variation
	vars, err := sweep.ParseVariations(*vary)
	if err != nil {
		fmt.Println("An error has occurred: ", err.Error())
		os.Exit(1)
	}

	// Define sweep options for calculations
	opt := sweep.SweepOptions{
		Runner: calc.RunnerOptions{
			EnableRTB:       *enableRTB,
			BitrateGbps:     *bitrate,
			TxRiseTimeNs:    *txrt,
			RxRiseTimeNs:    *rxrt,
			DispersionPerKm: *dispersion,
		},
	}
	results := sweep.RunSweep(links, []sweep.Variation{vars}, opt)
	if err := foio.WriteCSV(*output, results, ','); err != nil {
		fmt.Println("An error has occurred: ", err.Error())
		os.Exit(1)
	}
	fmt.Printf("DONE — %d swept links written to %s\n", len(results), *output)
}