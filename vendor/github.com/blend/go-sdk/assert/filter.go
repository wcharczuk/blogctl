package assert

import (
	"flag"
	"testing"
)

var (
	unit        = flag.Bool("unit", false, "If we should run unit tests")
	acceptance  = flag.Bool("acceptance", false, "If we should run acceptance tests")
	integration = flag.Bool("integration", false, "If we should run integration tests")
)

func init() {
	flag.Parse()
}

// Filter is a unit test filter.
type Filter string

// Filters
const (
	// Unit is a filter for unit tests.
	Unit = "unit"
	// Acceptance is a filter for acceptance tests.
	Acceptance = "acceptance"
	// Integration is a filter for integration tests.
	Integration = "integration"
)

// CheckFilter checks the filter.
func CheckFilter(t *testing.T, filter Filter) {
	if !*unit && !*acceptance && !*integration {
		return
	}

	switch filter {
	case Unit:
		if unit != nil && !*unit {
			t.Skip()
		}
	case Acceptance:
		if acceptance != nil && !*acceptance {
			t.Skip()
		}
	case Integration:
		if integration != nil && !*integration {
			t.Skip()
		}
	}
}
