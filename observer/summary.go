package observer

import (
	"fmt"
	"sort"

	"github.com/imkira/go-observer"
	"github.com/magnusbaeck/logstash-filter-verifier/logging"
)

// Summary summarizes the number of successful and failed test cases.
type Summary struct {
	NumberOk    int
	NumberNotOk int
}

var log = logging.MustGetLogger()

// RunSummaryObserver launches a consumer responsible for printing a summary
// at the end of the execution.
func RunSummaryObserver(prop observer.Property) {
	var (
		results       map[string]Summary
		globalSummary Summary
	)

	stream := prop.Observe()

	for {
		data := stream.Value()

		switch event := data.(type) {
		// Init struct to store result test
		case TestExecutionStart:
			results = make(map[string]Summary)
			globalSummary = Summary{
				NumberOk:    0,
				NumberNotOk: 0,
			}
		// Display result on stdout
		case TestExecutionEnd:
			fmt.Printf("\nSummary: %s All tests : %d/%d\n", getIconStatus(globalSummary.NumberNotOk), globalSummary.NumberOk, globalSummary.NumberOk+globalSummary.NumberNotOk)

			// Ordering by keys name
			keys := make([]string, len(results))
			i := 0
			for key := range results {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for _, key := range keys {
				summary := results[key]

				fmt.Printf("\t %s %s: %d/%d\n", getIconStatus(summary.NumberNotOk), key, summary.NumberOk, summary.NumberOk+summary.NumberNotOk)
			}
		// Store result test
		case ComparisonResult:

			// Compute summary to display at the end and display current test status
			summary := results[event.Path]
			if event.Status {
				summary.NumberOk++
				globalSummary.NumberOk++
				fmt.Printf("\u2611 %s from %s\n", event.Name, event.Path)
			} else {
				summary.NumberNotOk++
				globalSummary.NumberNotOk++
				fmt.Printf("\u2610 %s from %s:\n%s\n", event.Name, event.Path, event.Explain)
			}
			results[event.Path] = summary
		default:
			log.Debugf("Receive data that we doesn't say how to manage it %+v", data)
		}

		<-stream.Changes()
		stream.Next()
	}
}

func getIconStatus(numberNotOk int) string {
	if numberNotOk == 0 {
		return "\u2611"
	}

	return "\u2610"
}
